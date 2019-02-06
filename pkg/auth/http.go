package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/tracing/opentracing"

	"github.com/gorilla/mux"

	"github.com/cedric-parisi/payment-api/pkg/errorhandling"
	"github.com/cedric-parisi/payment-api/pkg/instrumenting"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type authRequest struct {
	ID string `json:"id"`
}

type authResponse struct {
	Token string `json:"token"`
}

// MakeAuthHandler ...
func MakeAuthHandler(service Service, errLogger kitlog.Logger, tracer stdopentracing.Tracer) http.Handler {
	errLogger = kitlog.With(errLogger, "component", "auth")

	options := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorEncoder(errorEncoder(errLogger)),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authRequest)
		token, err := service.GetJWT(req.ID)
		if err != nil {
			return nil, err
		}
		return authResponse{Token: token}, nil
	}

	authHandler := instrumenting.Middleware("auth", "post-auth",
		kithttp.NewServer(
			endpoint,
			decodeAuthRequest,
			kithttp.EncodeJSONResponse,
			append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "auth", errLogger)))...,
		),
	)
	r := mux.NewRouter().PathPrefix("/auth/").Subrouter().StrictSlash(true)
	{
		r.Handle("/", errorhandling.RecoverFromPanic(errLogger, authHandler)).Methods(http.MethodPost)
	}

	return r
}

func decodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request authRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func errorEncoder(logger kitlog.Logger) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		kithttp.DefaultErrorEncoder(ctx, err, w)
	}
}
