// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http 0.1.0

package v1

import (
	context "context"
	json "encoding/json"
	go_restful "github.com/emicklei/go-restful"
	errors "github.com/tkeel-io/kit/errors"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	reflect "reflect"
)

import transportHTTP "github.com/tkeel-io/kit/transport/http"

// This is a compile-time assertion to ensure that this generated file
// is compatible with the tkeel package it is being compiled against.
// import package.context.http.reflect.go_restful.json.errors.emptypb.

type EntityHTTPServer interface {
	AppendMapper(context.Context, *AppendMapperRequest) (*EntityResponse, error)
	CreateEntity(context.Context, *CreateEntityRequest) (*EntityResponse, error)
	DeleteEntity(context.Context, *DeleteEntityRequest) (*DeleteEntityResponse, error)
	GetEntity(context.Context, *GetEntityRequest) (*EntityResponse, error)
	ListEntity(context.Context, *ListEntityRequest) (*ListEntityResponse, error)
	PatchEntity(context.Context, *PatchEntityRequest) (*EntityResponse, error)
	PatchEntityZ(context.Context, *PatchEntityRequest) (*EntityResponse, error)
	SetEntityConfigs(context.Context, *SetEntityConfigRequest) (*EntityResponse, error)
	UpdateEntity(context.Context, *UpdateEntityRequest) (*EntityResponse, error)
}

type EntityHTTPHandler struct {
	srv EntityHTTPServer
}

func newEntityHTTPHandler(s EntityHTTPServer) *EntityHTTPHandler {
	return &EntityHTTPHandler{srv: s}
}

func (h *EntityHTTPHandler) AppendMapper(req *go_restful.Request, resp *go_restful.Response) {
	in := AppendMapperRequest{}
	if err := transportHTTP.GetBody(req, &in.Mapper); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.AppendMapper(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) CreateEntity(req *go_restful.Request, resp *go_restful.Response) {
	in := CreateEntityRequest{}
	if err := transportHTTP.GetBody(req, &in.Properties); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.CreateEntity(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) DeleteEntity(req *go_restful.Request, resp *go_restful.Response) {
	in := DeleteEntityRequest{}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.DeleteEntity(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) GetEntity(req *go_restful.Request, resp *go_restful.Response) {
	in := GetEntityRequest{}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.GetEntity(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) ListEntity(req *go_restful.Request, resp *go_restful.Response) {
	in := ListEntityRequest{}
	if err := transportHTTP.GetBody(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.ListEntity(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) PatchEntity(req *go_restful.Request, resp *go_restful.Response) {
	in := PatchEntityRequest{}
	if err := transportHTTP.GetBody(req, &in.Properties); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.PatchEntity(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) PatchEntityZ(req *go_restful.Request, resp *go_restful.Response) {
	in := PatchEntityRequest{}
	if err := transportHTTP.GetBody(req, &in.Properties); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.PatchEntityZ(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) SetEntityConfigs(req *go_restful.Request, resp *go_restful.Response) {
	in := SetEntityConfigRequest{}
	if err := transportHTTP.GetBody(req, &in.Configs); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.SetEntityConfigs(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *EntityHTTPHandler) UpdateEntity(req *go_restful.Request, resp *go_restful.Response) {
	in := UpdateEntityRequest{}
	if err := transportHTTP.GetBody(req, &in.Properties); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ctx := transportHTTP.ContextWithHeader(req.Request.Context(), req.Request.Header)

	out, err := h.srv.UpdateEntity(ctx, &in)
	if err != nil {
		tErr := errors.FromError(err)
		httpCode := errors.GRPCToHTTPStatusCode(tErr.GRPCStatus().Code())
		resp.WriteErrorString(httpCode, tErr.Message)
		return
	}
	if reflect.ValueOf(out).Elem().Type().AssignableTo(reflect.TypeOf(emptypb.Empty{})) {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func RegisterEntityHTTPServer(container *go_restful.Container, srv EntityHTTPServer) {
	var ws *go_restful.WebService
	for _, v := range container.RegisteredWebServices() {
		if v.RootPath() == "/v1" {
			ws = v
			break
		}
	}
	if ws == nil {
		ws = new(go_restful.WebService)
		ws.ApiVersion("/v1")
		ws.Path("/v1").Produces(go_restful.MIME_JSON)
		container.Add(ws)
	}

	handler := newEntityHTTPHandler(srv)
	ws.Route(ws.POST("/entities").
		To(handler.CreateEntity))
	ws.Route(ws.PUT("/entities/{id}").
		To(handler.UpdateEntity))
	ws.Route(ws.PATCH("/entities/{id}").
		To(handler.PatchEntity))
	ws.Route(ws.PUT("/entities/{id}/patch").
		To(handler.PatchEntityZ))
	ws.Route(ws.DELETE("/entities/{id}").
		To(handler.DeleteEntity))
	ws.Route(ws.GET("/entities/{id}").
		To(handler.GetEntity))
	ws.Route(ws.POST("/entities/search").
		To(handler.ListEntity))
	ws.Route(ws.POST("/entities/{id}/mappers").
		To(handler.AppendMapper))
	ws.Route(ws.PUT("/entities/{id}/configs").
		To(handler.SetEntityConfigs))
}
