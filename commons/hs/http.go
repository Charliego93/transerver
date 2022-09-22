package hs

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/transerver/commons"
)

func NewHTTPServerWithoutOpts(
	services []commons.Service,
) (*http.Server, error) {
	return NewHTTPServerWithOptions(services, nil)
}

func NewHTTPServerWithoutMuxOpts(
	services []commons.Service,
	handlers []Handler,
) (*http.Server, error) {
	return NewHTTPServerWithOptions(services, handlers)
}

func NewHTTPServerWithOptions(
	services []commons.Service,
	handlers []Handler,
	muxOpts ...runtime.ServeMuxOption,
) (*http.Server, error) {
	muxOpts = append(muxOpts, runtime.WithMarshalerOption("application/json", NewJSONMarshaler()))
	mux := runtime.NewServeMux(muxOpts...)
	for _, handler := range handlers {
		if err := mux.HandlePath(handler.Method, handler.Path, handler.route); err != nil {
			return nil, err
		}
	}

	for _, service := range services {
		if err := service.RegisterHTTP(mux); err != nil {
			return nil, err
		}
	}
	return &http.Server{Handler: mux}, nil
}

func DefaultOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(DefaultErrorHandler),
	}
}