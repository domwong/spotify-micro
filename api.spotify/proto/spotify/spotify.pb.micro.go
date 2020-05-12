// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/spotify/spotify.proto

package go_micro_api_spotify

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Spotify service

func NewSpotifyEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Spotify service

type SpotifyService interface {
	RootRedirect(ctx context.Context, in *RedirectRequest, opts ...client.CallOption) (*RedirectResponse, error)
	Callback(ctx context.Context, in *CallbackRequest, opts ...client.CallOption) (*CallbackResponse, error)
	Save(ctx context.Context, in *SaveRequest, opts ...client.CallOption) (*SaveResponse, error)
}

type spotifyService struct {
	c    client.Client
	name string
}

func NewSpotifyService(name string, c client.Client) SpotifyService {
	return &spotifyService{
		c:    c,
		name: name,
	}
}

func (c *spotifyService) RootRedirect(ctx context.Context, in *RedirectRequest, opts ...client.CallOption) (*RedirectResponse, error) {
	req := c.c.NewRequest(c.name, "Spotify.RootRedirect", in)
	out := new(RedirectResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spotifyService) Callback(ctx context.Context, in *CallbackRequest, opts ...client.CallOption) (*CallbackResponse, error) {
	req := c.c.NewRequest(c.name, "Spotify.Callback", in)
	out := new(CallbackResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spotifyService) Save(ctx context.Context, in *SaveRequest, opts ...client.CallOption) (*SaveResponse, error) {
	req := c.c.NewRequest(c.name, "Spotify.Save", in)
	out := new(SaveResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Spotify service

type SpotifyHandler interface {
	RootRedirect(context.Context, *RedirectRequest, *RedirectResponse) error
	Callback(context.Context, *CallbackRequest, *CallbackResponse) error
	Save(context.Context, *SaveRequest, *SaveResponse) error
}

func RegisterSpotifyHandler(s server.Server, hdlr SpotifyHandler, opts ...server.HandlerOption) error {
	type spotify interface {
		RootRedirect(ctx context.Context, in *RedirectRequest, out *RedirectResponse) error
		Callback(ctx context.Context, in *CallbackRequest, out *CallbackResponse) error
		Save(ctx context.Context, in *SaveRequest, out *SaveResponse) error
	}
	type Spotify struct {
		spotify
	}
	h := &spotifyHandler{hdlr}
	return s.Handle(s.NewHandler(&Spotify{h}, opts...))
}

type spotifyHandler struct {
	SpotifyHandler
}

func (h *spotifyHandler) RootRedirect(ctx context.Context, in *RedirectRequest, out *RedirectResponse) error {
	return h.SpotifyHandler.RootRedirect(ctx, in, out)
}

func (h *spotifyHandler) Callback(ctx context.Context, in *CallbackRequest, out *CallbackResponse) error {
	return h.SpotifyHandler.Callback(ctx, in, out)
}

func (h *spotifyHandler) Save(ctx context.Context, in *SaveRequest, out *SaveResponse) error {
	return h.SpotifyHandler.Save(ctx, in, out)
}
