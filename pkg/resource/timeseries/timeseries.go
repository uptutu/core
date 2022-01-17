package timeseries

import (
	"context"

	"github.com/tkeel-io/core/pkg/resource"

	"github.com/pkg/errors"
)

var ErrUnregisteredEngine = errors.New("unregistered engine")

var registerMap = map[Engine]Generator{
	EngineInflux: newInflux,
	EngineNoop:   newNoop,
}

type Actuator interface {
	Init(resource.TimeSeriesMetadata) error
	Write(ctx context.Context, req *WriteRequest) *Response
	Query(ctx context.Context, req QueryRequest) *Response
}

type Generator func() Actuator

func New(name Engine) (Actuator, error) {
	if generator, has := registerMap[name]; has {
		return generator(), nil
	}
	return nil, ErrUnregisteredEngine
}

func Register(name Engine, handler Generator) {
	registerMap[name] = handler
}
