package server

import (
	sj "encoding/json"
	nt "net/http"
	"shop/api/shop/v1"
	"shop/internal/conf"
	"shop/internal/service"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, shop *service.ShopService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.ErrorEncoder(DefaultErrorEncoder),
		http.ResponseEncoder(DefaultResponseEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterShopHTTPServer(srv, shop)
	return srv
}

// DefaultResponseEncoder copy from http.DefaultResponseEncoder
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		nt.Redirect(w, r, url, code)
		return nil
	}

	codec := encoding.GetCodec(json.Name) // ignore Accept Header
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}

	bs, _ := sj.Marshal(sj.RawMessage(data))

	w.Header().Set("Content-Type", ContentType(codec.Name()))
	_, err = w.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

// DefaultErrorEncoder copy from http.DefaultErrorEncoder.
func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	errorCode := int(errors.FromError(err).Code)
	codec := encoding.GetCodec(json.Name) // ignore Accept header
	var body []byte
	var se *ErrorResponse
	if errors.Is(err, service.ErrBadRequest) || errors.Is(err, service.ErrUnauthorized) {
		se = FromError(errors.FromError(err)) // change error to BaseResponse
	} else {
		se = &ErrorResponse{Message: "Неизвестная ошибка"}
		errorCode = 500
	}
	body, err = codec.Marshal(se)
	if err != nil {
		w.WriteHeader(nt.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	w.WriteHeader(errorCode)
	// w.WriteHeader(int(se.Code)) // ignore http status code
	_, _ = w.Write(body)
}

const (
	baseContentType = "application"
)

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}

func NewResponse(data []byte) BaseResponse {
	return BaseResponse{
		Data:    sj.RawMessage(data),
	}
}

func FromError(e *errors.Error) *ErrorResponse {
	if e == nil {
		return nil
	}
	return &ErrorResponse{
		Message: e.Message,
	}
}

type BaseResponse struct {
	Data    sj.RawMessage `"omitempty"`
}

type ErrorResponse struct {
	Message string        `json:"errors"`
}
