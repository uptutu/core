/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runtime

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	ants "github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	pb "github.com/tkeel-io/core/api/core/v1"
	"github.com/tkeel-io/core/pkg/config"
	"github.com/tkeel-io/core/pkg/constraint"
	"github.com/tkeel-io/core/pkg/logger"
	"github.com/tkeel-io/core/pkg/resource"
	"github.com/tkeel-io/core/pkg/resource/tseries"
	"github.com/tkeel-io/core/pkg/statem"
	"github.com/tkeel-io/core/pkg/util"
	"github.com/tkeel-io/kit/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type Manager struct {
	containers    map[string]*Container
	msgCh         chan statem.MessageContext
	disposeCh     chan statem.MessageContext
	coroutinePool *ants.Pool
	actorEnv      *Environment

	daprClient    dapr.Client
	etcdClient    *clientv3.Client
	searchClient  pb.SearchHTTPServer
	tseriesClient tseries.TimeSerier

	shutdown chan struct{}
	lock     sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewManager(ctx context.Context, coroutinePool *ants.Pool, searchClient pb.SearchHTTPServer) (*Manager, error) {
	var (
		err        error
		daprClient dapr.Client
		etcdClient *clientv3.Client
	)

	expireTime := 3 * time.Second
	etcdAddr := config.GetConfig().Etcd.Address
	tseriesClient := tseries.NewTimeSerier(config.GetConfig().TimeSeries.Name)
	returnErr := func(err error) error { return errors.Wrap(err, "create manager failed") }

	if daprClient, err = dapr.NewClient(); nil != err {
		return nil, returnErr(err)
	} else if err = tseriesClient.Init(resource.ParseFrom(&config.GetConfig().TimeSeries)); nil != err {
		return nil, returnErr(err)
	} else if etcdClient, err = clientv3.New(clientv3.Config{Endpoints: etcdAddr, DialTimeout: expireTime}); nil != err {
		return nil, returnErr(err)
	}

	ctx, cancel := context.WithCancel(ctx)

	mgr := &Manager{
		ctx:           ctx,
		cancel:        cancel,
		actorEnv:      NewEnv(),
		daprClient:    daprClient,
		etcdClient:    etcdClient,
		searchClient:  searchClient,
		tseriesClient: tseriesClient,
		containers:    make(map[string]*Container),
		msgCh:         make(chan statem.MessageContext, 10),
		disposeCh:     make(chan statem.MessageContext, 10),
		coroutinePool: coroutinePool,
		lock:          sync.RWMutex{},
	}

	// set default container.
	mgr.containers["default"] = NewContainer()
	return mgr, nil
}

func (m *Manager) SendMsg(msgCtx statem.MessageContext) {
	bytes, _ := json.Marshal(msgCtx)
	log.Debug("actor send message", zap.String("msg", string(bytes)))

	// 解耦actor之间的直接调用
	m.msgCh <- msgCtx
}

func (m *Manager) init() error {
	// load all subcriptions.
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	log.Info("initialize actor manager, tql loadding...")
	res, err := m.etcdClient.Get(ctx, TQLEtcdPrefix, clientv3.WithPrefix())
	if nil != err {
		return errors.Wrap(err, "load all tql")
	}

	descs := make([]EtcdPair, len(res.Kvs))
	for index, kv := range res.Kvs {
		descs[index] = EtcdPair{Key: string(kv.Key), Value: kv.Value}
	}

	stateIDs := make([]string, 0)
	loadEntities := m.actorEnv.LoadMapper(descs)
	for _, info := range loadEntities {
		log.Debug("load state marchine", logger.EntityID(info.EntityID), zap.String("type", info.Type))
		if err = m.loadActor(context.Background(), info.Type, info.EntityID); nil != err {
			log.Error("load state marchine", zap.Error(err), logger.EntityID(info.EntityID), zap.String("type", info.Type))
		}
		stateIDs = append(stateIDs, info.EntityID)
	}
	m.reloadActor(stateIDs)
	return nil
}

func (m *Manager) watchResource() error {
	// watch tqls.
	tqlWatcher, err := util.NewWatcher(m.ctx, config.GetConfig().Etcd.Address)
	if nil != err {
		return errors.Wrap(err, "create tql watcher failed")
	}

	log.Info("watch resource")
	tqlWatcher.Watch(TQLEtcdPrefix, true, func(ev *clientv3.Event) {
		// on changed.
		effects, _ := m.actorEnv.OnMapperChanged(ev.Type, EtcdPair{Key: string(ev.Kv.Key), Value: ev.Kv.Value})
		m.reloadActor(effects)
	})

	return nil
}

func (m *Manager) isThisNode() bool {
	return true
}

func (m *Manager) reloadActor(stateIDs []string) error {
	// 判断 actor 是否在当前节点.
	if m.isThisNode() {
		var err error
		for _, stateID := range stateIDs {
			var stateMarchine statem.StateMarchiner
			base := &statem.Base{ID: stateID, Type: StateMarchineTypeBasic}
			if _, stateMarchine = m.getStateMarchine("", stateID); nil != stateMarchine {
				log.Debug("load state marchine @ runtime.", logger.EntityID(stateID))
			} else if stateMarchine, err = m.loadOrCreate(m.ctx, "", false, base); nil == err {
				stateMarchine.LoadEnvironments(m.actorEnv.GetEnvBy(stateID))
				continue
			}
		}
	}
	return nil
}

func (m *Manager) Start() error {
	// init: load some resource.
	m.init()
	// watch resource.
	m.watchResource()

	go func() {
		for {
			select {
			case <-m.ctx.Done():
				log.Info("entity manager exited.")
				return
			case msgCtx := <-m.msgCh:
				// dispatch message. 将消息分发到不同的节点。
				m.disposeCh <- msgCtx

			case msgCtx := <-m.disposeCh:
				eid := msgCtx.Headers.GetTargetID()
				channelID := msgCtx.Headers.Get(statem.MessageCtxHeaderChannelID)
				log.Debug("dispose message", logger.EntityID(eid), logger.MessageInst(msgCtx))
				channelID, stateMarchine := m.getStateMarchine(channelID, eid)
				if nil == stateMarchine {
					var err error
					en := &statem.Base{
						ID:     eid,
						Owner:  msgCtx.Headers.GetOwner(),
						Source: msgCtx.Headers.GetSource(),
						Type:   msgCtx.Headers.Get(statem.MessageCtxHeaderType),
					}
					stateMarchine, err = m.loadOrCreate(m.ctx, channelID, true, en)
					if nil != err {
						log.Error("dispatching message", zap.Error(err),
							logger.EntityID(eid), zap.String("channel", channelID), logger.MessageInst(msgCtx))
						continue
					}
				}

				if stateMarchine.OnMessage(msgCtx.Message) {
					// attatch goroutine to entity.
					m.coroutinePool.Submit(stateMarchine.HandleLoop)
				}
			case <-m.shutdown:
				log.Info("state marchine manager exit.")
				return
			}
		}
	}()

	return nil
}

func (m *Manager) Shutdown() {
	m.cancel()
	m.shutdown <- struct{}{}
}

func (m *Manager) GetDaprClient() dapr.Client {
	return m.daprClient
}

func (m *Manager) getStateMarchine(cid, eid string) (string, statem.StateMarchiner) { //nolint
	if cid == "" {
		cid = "default"
	}

	if container, ok := m.containers[cid]; ok {
		if sm := container.Get(eid); nil != sm {
			if sm.GetStatus() == statem.SMStatusDeleted {
				container.Remove(eid)
				return cid, nil
			}
			return cid, sm
		}
	}

	for channelID, container := range m.containers {
		if sm := container.Get(eid); sm != nil {
			if sm.GetStatus() == statem.SMStatusDeleted {
				container.Remove(eid)
				return cid, nil
			}

			if channelID == "default" && cid != channelID {
				container.Remove(sm.GetID())
				if _, ok := m.containers[cid]; !ok {
					m.containers[cid] = NewContainer()
				}
				m.containers[cid].Add(sm)
			}
			return cid, sm
		}
	}

	return cid, nil
}

func (m *Manager) loadActor(ctx context.Context, typ string, id string) error {
	_, err := m.loadOrCreate(ctx, "", false, &statem.Base{ID: id, Type: typ})
	return errors.Wrap(err, "load entity")
}

func (m *Manager) loadOrCreate(ctx context.Context, channelID string, flagCreate bool, base *statem.Base) (sm statem.StateMarchiner, err error) { // nolint
	var en *statem.Base
	var res *dapr.StateItem
	res, err = m.daprClient.GetState(ctx, EntityStateName, base.ID)

	if nil != err && !flagCreate {
		return nil, errors.Wrap(err, "load state marchine")
	} else if en, err = statem.DecodeBase(res.Value); nil == err {
		base = en // decode value to statem.Base.
	} else if !flagCreate {
		return nil, errors.Wrap(err, "load state marchine, state not found")
	}

	log.Debug("load or create state marchiner",
		logger.EntityID(base.ID),
		zap.String("type", base.Type),
		zap.String("owner", base.Owner),
		zap.String("source", base.Source))

	switch base.Type {
	case StateMarchineTypeSubscription:
		if sm, err = newSubscription(ctx, m, base); nil != err {
			return nil, errors.Wrap(err, "load subscription")
		}
	default:
		// default base entity type.
		if sm, err = statem.NewState(ctx, m, base, nil); nil != err {
			return nil, errors.Wrap(err, "load state marchine")
		}
	}

	if channelID == "" {
		channelID = "defult"
	}

	if _, has := m.containers[channelID]; !has {
		m.containers[channelID] = NewContainer()
	}

	sm.Setup()
	m.containers[channelID].Add(sm)
	return sm, nil
}

func (m *Manager) HandleMsg(ctx context.Context, msg statem.MessageContext) {
	// dispose message from pubsub.
	m.msgCh <- msg
}

// Tools.

func (m *Manager) EscapedEntities(expression string) []string {
	return []string{expression}
}

// ------------------------------------APIs-----------------------------.

// SetProperties set properties into entity.
func (m *Manager) SetProperties(ctx context.Context, en *statem.Base) error {
	if en.ID == "" {
		en.ID = uuid()
	}

	// set properties.
	msgCtx := statem.MessageContext{
		Headers: statem.Header{},
		Message: statem.PropertyMessage{
			StateID:    en.ID,
			Operator:   constraint.PatchOpReplace.String(),
			Properties: en.KValues,
		},
	}
	msgCtx.Headers.SetOwner(en.Owner)
	msgCtx.Headers.SetTargetID(en.ID)
	msgCtx.Headers.SetSource(en.Source)
	msgCtx.Headers.Set(statem.MessageCtxHeaderType, en.Type)

	m.SendMsg(msgCtx)

	return nil
}

func (m *Manager) PatchEntity(ctx context.Context, en *statem.Base, patchData []*pb.PatchData) error {
	pdm := make(map[string][]*pb.PatchData)
	for _, pd := range patchData {
		pdm[pd.Operator] = append(pdm[pd.Operator], pd)
	}

	for op, pds := range pdm {
		kvs := make(map[string]constraint.Node)
		for _, pd := range pds {
			kvs[pd.Path] = constraint.NewNode(pd.Value.AsInterface())
		}

		if len(kvs) > 0 {
			msgCtx := statem.MessageContext{
				Headers: statem.Header{},
				Message: statem.PropertyMessage{
					StateID:    en.ID,
					Operator:   op,
					Properties: kvs,
				},
			}

			// set headers.
			msgCtx.Headers.SetOwner(en.Owner)
			msgCtx.Headers.SetTargetID(en.ID)
			msgCtx.Headers.Set(statem.MessageCtxHeaderType, en.Type)
			m.SendMsg(msgCtx)
		}
	}

	return nil
}

// SetConfigs set entity configs.
func (m *Manager) SetConfigs(ctx context.Context, en *statem.Base) error {
	var (
		err           error
		channelID     string
		stateMarchine statem.StateMarchiner
	)

	// load state marchine.
	if channelID, stateMarchine = m.getStateMarchine("", en.ID); nil == stateMarchine {
		if stateMarchine, err = m.loadOrCreate(ctx, channelID, false, en); nil != err {
			log.Error("set configs",
				logger.EntityID(en.ID),
				zap.Any("entity", en),
				zap.String("channel", channelID))
			return errors.Wrap(err, "set entity configs")
		}
	}

	// set entity configs.
	if err = stateMarchine.SetConfigs(en.Configs); nil != err {
		return errors.Wrap(err, "set entity configs")
	}

	// flush entity configs.
	return errors.Wrap(stateMarchine.Flush(ctx), "set entity configs")
}

// PatchConfigs patch entity configs.
func (m *Manager) PatchConfigs(ctx context.Context, en *statem.Base, patchData []*statem.PatchData) error {
	var (
		err           error
		channelID     string
		stateMarchine statem.StateMarchiner
	)

	// load state marchine.
	if channelID, stateMarchine = m.getStateMarchine("", en.ID); nil == stateMarchine {
		if stateMarchine, err = m.loadOrCreate(ctx, channelID, false, en); nil != err {
			log.Error("set configs",
				logger.EntityID(en.ID),
				zap.Any("entity", en),
				zap.String("channel", channelID))
			return errors.Wrap(err, "set entity configs")
		}
	}

	// set entity configs.
	if err = stateMarchine.PatchConfigs(patchData); nil != err {
		return errors.Wrap(err, "set entity configs")
	}

	// flush entity configs.
	return errors.Wrap(stateMarchine.Flush(ctx), "set entity configs")
}

// AppendConfigs append entity configs.
func (m *Manager) AppendConfigs(ctx context.Context, en *statem.Base) error {
	var (
		err           error
		channelID     string
		stateMarchine statem.StateMarchiner
	)

	// load state marchine.
	if channelID, stateMarchine = m.getStateMarchine("", en.ID); nil == stateMarchine {
		if stateMarchine, err = m.loadOrCreate(ctx, channelID, false, en); nil != err {
			log.Error("append configs",
				logger.EntityID(en.ID),
				zap.Any("entity", en),
				zap.String("channel", channelID))
			return errors.Wrap(err, "append entity configs")
		}
	}

	// append entity configs.
	if err = stateMarchine.AppendConfigs(en.Configs); nil != err {
		return errors.Wrap(err, "append entity configs")
	}

	// flush entity configs.
	return errors.Wrap(stateMarchine.Flush(ctx), "append entity configs")
}

// RemoveConfigs remove entity configs.
func (m *Manager) RemoveConfigs(ctx context.Context, en *statem.Base, propertyIDs []string) error {
	var (
		err           error
		channelID     string
		stateMarchine statem.StateMarchiner
	)

	// load state marchine.
	if channelID, stateMarchine = m.getStateMarchine("", en.ID); nil == stateMarchine {
		if stateMarchine, err = m.loadOrCreate(ctx, channelID, false, en); nil != err {
			log.Error("remove configs",
				logger.EntityID(en.ID),
				zap.Any("entity", en),
				zap.String("channel", channelID))
			return errors.Wrap(err, "remove entity configs")
		}
	}

	// remove entity configs.
	if err = stateMarchine.RemoveConfigs(propertyIDs); nil != err {
		return errors.Wrap(err, "remove entity configs")
	}

	// flush entity configs.
	return errors.Wrap(stateMarchine.Flush(ctx), "remove entity configs")
}

// DeleteStateMarchin delete runtime.Entity.
func (m *Manager) DeleteStateMarchin(ctx context.Context, base *statem.Base) (*statem.Base, error) {
	var err error
	channelID, stateMarchine := m.getStateMarchine("", base.ID)
	if nil == stateMarchine {
		if stateMarchine, err = m.loadOrCreate(m.ctx, channelID, true, base); nil != err {
			log.Error("remove configs",
				logger.EntityID(base.ID),
				zap.Any("entity", base),
				zap.String("channel", channelID))
			return nil, errors.Wrap(err, "remove entity configs")
		}
	}
	stateMarchine.SetStatus(statem.SMStatusDeleted)
	return stateMarchine.GetBase(), nil
}

// CleanEntity clean entity.
func (m *Manager) CleanEntity(ctx context.Context, id string) error {
	channelID, sm := m.getStateMarchine("", id)
	if nil != sm {
		m.containers[channelID].Remove(id)
	}
	return nil
}

// AppendMapper append a mapper into entity.
func (m *Manager) AppendMapper(ctx context.Context, en *statem.Base) error {
	if len(en.Mappers) == 0 {
		log.Error("append mapper into entity failed.", logger.EntityID(en.ID), zap.Error(ErrInvalidParams))
		return errors.Wrap(ErrInvalidParams, "append entity mapper failed")
	}

	msgCtx := statem.MessageContext{
		Headers: statem.Header{},
		Message: statem.MapperMessage{
			Operator: statem.MapperOperatorAppend,
			Mapper:   en.Mappers[0],
		},
	}

	msgCtx.Headers.SetOwner(en.Owner)
	msgCtx.Headers.SetTargetID(en.ID)

	m.SendMsg(msgCtx)

	return nil
}

// DeleteMapper delete mapper from entity.
func (m *Manager) RemoveMapper(ctx context.Context, en *statem.Base) error {
	if len(en.Mappers) == 0 {
		log.Error("remove mapper failed.", logger.EntityID(en.ID), zap.Error(ErrInvalidParams))
		return errors.Wrap(ErrInvalidParams, "remove entity mapper failed")
	}

	msgCtx := statem.MessageContext{
		Headers: statem.Header{},
		Message: statem.MapperMessage{
			Operator: statem.MapperOperatorRemove,
			Mapper:   en.Mappers[0],
		},
	}

	msgCtx.Headers.SetOwner(en.Owner)
	msgCtx.Headers.SetTargetID(en.ID)

	m.SendMsg(msgCtx)

	return nil
}

func (m *Manager) SearchFlush(ctx context.Context, values map[string]interface{}) error {
	var err error
	var val *structpb.Value
	if val, err = structpb.NewValue(values); nil != err {
		log.Error("search index failed.", zap.Error(err))
	} else if _, err = m.searchClient.Index(ctx, &pb.IndexObject{Obj: val}); nil != err {
		log.Error("search index failed.", zap.Error(err))
	}
	return errors.Wrap(err, "SearchFlushfailed")
}

func (m *Manager) TimeSeriesFlush(ctx context.Context, tds []tseries.TSeriesData) error {
	var err error
	for _, data := range tds {
		data.Fields["value"] = data.Value
		line := fmt.Sprintf("%s,%s %s", data.Measurement, util.ExtractMap(data.Tags), util.ExtractMap(data.Fields))

		_, err = m.tseriesClient.Write(ctx, &tseries.TSeriesRequest{
			Data:     []string{line},
			Metadata: map[string]string{},
		})

		if nil != err {
			// 这里其实是有问题的, 如果没写成功，怎么处理，和MQ的ack相关，考虑放到batch_queue处理.
			log.Error("flush time series data failed", zap.Error(err), zap.Any("data", data))
		}
	}

	return errors.Wrap(err, "TimeSeriesFlush")
}

// uuid generate an uuid.
func uuid() string {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return ""
	}
	// see section 4.1.1.
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// see section 4.1.3.
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
