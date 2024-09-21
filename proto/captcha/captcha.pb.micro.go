// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/captcha.proto

package captcha

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Captcha service

func NewCaptchaEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Captcha service

type CaptchaService interface {
	GenerateCaptcha(ctx context.Context, in *GenerateCaptchaRequest, opts ...client.CallOption) (*GenerateCaptchaResponse, error)
	VerifyCaptcha(ctx context.Context, in *VerifyCaptchaRequest, opts ...client.CallOption) (*VerifyCaptchaResponse, error)
}

type captchaService struct {
	c    client.Client
	name string
}

func NewCaptchaService(name string, c client.Client) CaptchaService {
	return &captchaService{
		c:    c,
		name: name,
	}
}

func (c *captchaService) GenerateCaptcha(ctx context.Context, in *GenerateCaptchaRequest, opts ...client.CallOption) (*GenerateCaptchaResponse, error) {
	req := c.c.NewRequest(c.name, "Captcha.GenerateCaptcha", in)
	out := new(GenerateCaptchaResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *captchaService) VerifyCaptcha(ctx context.Context, in *VerifyCaptchaRequest, opts ...client.CallOption) (*VerifyCaptchaResponse, error) {
	req := c.c.NewRequest(c.name, "Captcha.VerifyCaptcha", in)
	out := new(VerifyCaptchaResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Captcha service

type CaptchaHandler interface {
	GenerateCaptcha(context.Context, *GenerateCaptchaRequest, *GenerateCaptchaResponse) error
	VerifyCaptcha(context.Context, *VerifyCaptchaRequest, *VerifyCaptchaResponse) error
}

func RegisterCaptchaHandler(s server.Server, hdlr CaptchaHandler, opts ...server.HandlerOption) error {
	type captcha interface {
		GenerateCaptcha(ctx context.Context, in *GenerateCaptchaRequest, out *GenerateCaptchaResponse) error
		VerifyCaptcha(ctx context.Context, in *VerifyCaptchaRequest, out *VerifyCaptchaResponse) error
	}
	type Captcha struct {
		captcha
	}
	h := &captchaHandler{hdlr}
	return s.Handle(s.NewHandler(&Captcha{h}, opts...))
}

type captchaHandler struct {
	CaptchaHandler
}

func (h *captchaHandler) GenerateCaptcha(ctx context.Context, in *GenerateCaptchaRequest, out *GenerateCaptchaResponse) error {
	return h.CaptchaHandler.GenerateCaptcha(ctx, in, out)
}

func (h *captchaHandler) VerifyCaptcha(ctx context.Context, in *VerifyCaptchaRequest, out *VerifyCaptchaResponse) error {
	return h.CaptchaHandler.VerifyCaptcha(ctx, in, out)
}