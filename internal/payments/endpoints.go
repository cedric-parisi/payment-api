package payments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cedric-parisi/payment-api/internal/models"
	"github.com/cedric-parisi/payment-api/pkg/utils"

	"github.com/go-kit/kit/endpoint"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// Endpoints ...
type Endpoints struct {
	CreatePayment       endpoint.Endpoint
	UpdatePayment       endpoint.Endpoint
	GetPayment          endpoint.Endpoint
	GetFilteredPayments endpoint.Endpoint
	DeletePayment       endpoint.Endpoint
}

// MakeEndpoints create endpoits
// With tracing and auth middleware
func MakeEndpoints(service Service, tracer opentracing.Tracer, JWTMiddleware endpoint.Middleware) Endpoints {
	return Endpoints{
		CreatePayment:       kitopentracing.TraceServer(tracer, "create_payment")(JWTMiddleware(MakeCreatePaymentEndpoint(service))),
		UpdatePayment:       kitopentracing.TraceServer(tracer, "update_payment")(JWTMiddleware(MakeUpdatePaymentEndpoint(service))),
		GetPayment:          kitopentracing.TraceServer(tracer, "get_payment")(MakeGetPaymentEndpoint(service)),
		GetFilteredPayments: kitopentracing.TraceServer(tracer, "get_filtered-payments")(MakeGetFilteredPaymentsEndpoint(service)),
		DeletePayment:       kitopentracing.TraceServer(tracer, "delete_payment")(JWTMiddleware(MakeDeletePaymentEndpoint(service))),
	}
}

// CreatePaymentResponse represents the response body for a payment creation request
// Contains the newly created payment
type CreatePaymentResponse struct {
	models.Payment
}

// StatusCode will set 201 for payment creation
func (c CreatePaymentResponse) StatusCode() int {
	return http.StatusCreated
}

// Headers will set Location header with the location of the newly created payment resource
func (c CreatePaymentResponse) Headers() http.Header {
	return http.Header{
		"Location": []string{fmt.Sprintf("/payments/%s", c.ID)},
	}
}

// MakeCreatePaymentEndpoint ...
func MakeCreatePaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Payment)
		res, err := s.CreatePayment(ctx, req)
		if err != nil {
			return nil, err
		}
		return CreatePaymentResponse{
			Payment: *res,
		}, nil
	}
}

// MakeUpdatePaymentEndpoint ...
func MakeUpdatePaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Payment)
		err := s.UpdatePayment(ctx, req)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// MakeGetPaymentEndpoint ...
func MakeGetPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		return s.GetPayment(ctx, id)
	}
}

// MakeGetFilteredPaymentsEndpoint ...
func MakeGetFilteredPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		filters := request.(*utils.Filter)
		return s.GetFilteredPayments(ctx, filters)
	}
}

// MakeDeletePaymentEndpoint ...
func MakeDeletePaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		err := s.DeletePayment(ctx, id)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
