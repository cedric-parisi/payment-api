package payments

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cedric-parisi/payment-api/internal/models"

	"github.com/cedric-parisi/payment-api/pkg/utils"

	"github.com/go-kit/kit/tracing/opentracing"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/cedric-parisi/payment-api/pkg/errorhandling"
	"github.com/cedric-parisi/payment-api/pkg/instrumenting"
)

// MakePaymentHTTPHandler ...
func MakePaymentHTTPHandler(errLogger kitlog.Logger, tracer stdopentracing.Tracer, endpoints Endpoints) http.Handler {
	errLogger = kitlog.With(errLogger, "component", resourceName)

	options := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorEncoder(encodeError(errLogger)),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	createPaymentHandler := instrumenting.Middleware(resourceName, "create-payment",
		kithttp.NewServer(
			endpoints.CreatePayment,
			decodeCreatePaymentRequest,
			kithttp.EncodeJSONResponse,
			append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "payments", errLogger)))...,
		),
	)

	updatePaymentHandler := instrumenting.Middleware(resourceName, "update-payment",
		kithttp.NewServer(
			endpoints.UpdatePayment,
			decodeUpdatePaymentRequest,
			encodeEmptyResponse,
			append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "payments", errLogger)))...,
		),
	)

	getPaymentHandler := instrumenting.Middleware(resourceName, "get-payment-by-id",
		kithttp.NewServer(
			endpoints.GetPayment,
			decodeGetPaymentRequest,
			kithttp.EncodeJSONResponse,
			append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "payments", errLogger)))...,
		),
	)

	getFilteredPaymentsHandler := instrumenting.Middleware(resourceName, "get-filtered-payments",
		kithttp.NewServer(
			endpoints.GetFilteredPayments,
			decodeGetFilteredPaymentsRequest,
			kithttp.EncodeJSONResponse,
			append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "payments", errLogger)))...,
		),
	)

	deletePaymentHandler := instrumenting.Middleware(resourceName, "delete-payment",
		kithttp.NewServer(
			endpoints.DeletePayment,
			decodeDeletePaymentRequest,
			encodeEmptyResponse,
			append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "payments", errLogger)))...,
		),
	)

	r := mux.NewRouter().PathPrefix("/payments/").Subrouter().StrictSlash(true)
	{
		r.Handle("/", errorhandling.RecoverFromPanic(errLogger, createPaymentHandler)).Methods(http.MethodPost)
		r.Handle("/{id}", errorhandling.RecoverFromPanic(errLogger, updatePaymentHandler)).Methods(http.MethodPut)
		r.Handle("/{id}", errorhandling.RecoverFromPanic(errLogger, getPaymentHandler)).Methods(http.MethodGet)
		r.Handle("/", errorhandling.RecoverFromPanic(errLogger, getFilteredPaymentsHandler)).Methods(http.MethodGet)
		r.Handle("/{id}", errorhandling.RecoverFromPanic(errLogger, deletePaymentHandler)).Methods(http.MethodDelete)
	}

	return r
}

func decodeCreatePaymentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := &models.Payment{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errorhandling.Internal(invalidPaymentCode, err)
	}
	return req, nil
}

func decodeUpdatePaymentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	req := &models.Payment{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errorhandling.Internal(invalidPaymentCode, err)
	}
	if id != req.ID.String() {
		return nil, errorhandling.InvalidRequest(invalidPaymentCode, errors.New("id mismatch"))
	}
	return req, nil
}

func decodeGetPaymentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return mux.Vars(r)["id"], nil
}

func decodeGetFilteredPaymentsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	filter := utils.GetFilter(r.URL.Query())
	return filter, nil
}

func decodeDeletePaymentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return mux.Vars(r)["id"], nil
}

func encodeEmptyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// encodeError logs internal errors before calling the default error encoder
// And catches errors that are not implemeting the APIError interface
func encodeError(logger kitlog.Logger) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		// An error was raised by the authentication middleware
		if err == kitjwt.ErrTokenContextMissing ||
			err == kitjwt.ErrTokenExpired ||
			err == kitjwt.ErrTokenInvalid ||
			err == kitjwt.ErrTokenMalformed ||
			err == kitjwt.ErrTokenNotActive ||
			err == kitjwt.ErrUnexpectedSigningMethod {
			err = errorhandling.Unauthorized("invalid_authentication_token", err)
		}

		if tmp, ok := err.(kithttp.StatusCoder); ok {
			// We log only internal server error
			if tmp.StatusCode() == http.StatusInternalServerError {
				errorhandling.Log(ctx, err, logger)
			}
		} else {
			// An error occured that was not catched by our error handling
			err = errorhandling.Internal("unknown_error", err)
			errorhandling.Log(ctx, err, logger)
		}

		// Use JSON default encoder from go-kit
		kithttp.DefaultErrorEncoder(ctx, err, w)
	}
}
