// +build !integration

package payments

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/cedric-parisi/payment-api/pkg/errorhandling"

	"github.com/cedric-parisi/payment-api/internal/models"
	"github.com/cedric-parisi/payment-api/pkg/utils"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_decodeCreatePaymentRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "decode create payment request ok",
			args: args{
				ctx: context.Background(),
				r: httptest.NewRequest(http.MethodPost, "/payments/", strings.NewReader(`
				{
					"type": "Withdraw",
					"version": 0,
					"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"
				}
				`)),
			},
			want: &models.Payment{
				Type:           models.WithdrawType,
				Version:        0,
				OrganisationID: uuid.MustParse("743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"),
			},
		},
		{
			name: "decode create payment request failed due to no body",
			args: args{
				ctx: context.Background(),
				r:   httptest.NewRequest(http.MethodPost, "/payments/", strings.NewReader(``)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeCreatePaymentRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCreatePaymentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeCreatePaymentRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeUpdatePaymentRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "decode update payment request ok",
			args: args{
				ctx: context.Background(),
				r: mux.SetURLVars(httptest.NewRequest(http.MethodPut, "/payments/3578205f-aeb3-444a-a42f-d47298b6eb8b", strings.NewReader(`
				{
					"id": "3578205f-aeb3-444a-a42f-d47298b6eb8b",
					"type": "Payment",
					"version": 18,
					"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"
				}
				`)), map[string]string{"id": "3578205f-aeb3-444a-a42f-d47298b6eb8b"}),
			},
			want: &models.Payment{
				ID:             uuid.MustParse("3578205f-aeb3-444a-a42f-d47298b6eb8b"),
				Type:           models.PaymentType,
				Version:        18,
				OrganisationID: uuid.MustParse("743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"),
			},
		},
		{
			name: "decode update payment request failed due to no body",
			args: args{
				ctx: context.Background(),
				r: mux.SetURLVars(httptest.NewRequest(http.MethodPut, "/payments/3578205f-aeb3-444a-a42f-d47298b6eb8b", strings.NewReader(``)),
					map[string]string{"id": "3578205f-aeb3-444a-a42f-d47298b6eb8b"}),
			},
			wantErr: true,
		},
		{
			name: "decode update payment failed due to wrong id",
			args: args{
				ctx: context.Background(),
				r: mux.SetURLVars(httptest.NewRequest(http.MethodPut, "/payments/3578205f-aeb3-444a-a42f-d47298b6eb8b", strings.NewReader(`
				{
					"id": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
					"type": "Payment",
					"version": 18,
					"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"
				}
				`)), map[string]string{"id": "3578205f-aeb3-444a-a42f-d47298b6eb8b"}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeUpdatePaymentRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeUpdatePaymentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeUpdatePaymentRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeGetPaymentRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "get payment request ok",
			args: args{
				ctx: context.Background(),
				r: mux.SetURLVars(
					httptest.NewRequest(http.MethodGet, "/offers/3578205f-aeb3-444a-a42f-d47298b6eb8b", nil),
					map[string]string{"id": "3578205f-aeb3-444a-a42f-d47298b6eb8b"}),
			},
			want: "3578205f-aeb3-444a-a42f-d47298b6eb8b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeGetPaymentRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeGetPaymentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeGetPaymentRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeGetFilteredPaymentsRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "get filtered payment request ok",
			args: args{
				ctx: context.Background(),
				r:   httptest.NewRequest(http.MethodGet, "/payments/?offset=0&limit=10&sort=-version", nil),
			},
			want: &utils.Filter{
				Limit:  10,
				Offset: 0,
				Sorting: []utils.Sort{
					{
						Field:      "version",
						Descending: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeGetFilteredPaymentsRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeGetFilteredPaymentsRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeGetFilteredPaymentsRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeDeletePaymentRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "delete payment request ok",
			args: args{
				ctx: context.Background(),
				r: mux.SetURLVars(
					httptest.NewRequest(http.MethodDelete, "/offers/3578205f-aeb3-444a-a42f-d47298b6eb8b", nil),
					map[string]string{"id": "3578205f-aeb3-444a-a42f-d47298b6eb8b"}),
			},
			want: "3578205f-aeb3-444a-a42f-d47298b6eb8b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeDeletePaymentRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeDeletePaymentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeDeletePaymentRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeEmptyResponse(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()

	// Act
	encodeEmptyResponse(context.Background(), w, strings.NewReader("here's a body"))

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func Test_encodeError(t *testing.T) {
	var tests = []struct {
		name string
		in   error
		out  int
	}{
		{
			name: "unknown error catched",
			in:   errors.New("not handled error"),
			out:  http.StatusInternalServerError,
		},
		{
			name: "jwt error catched",
			in:   kitjwt.ErrTokenInvalid,
			out:  http.StatusUnauthorized,
		},
		{
			name: "internal error",
			in:   errorhandling.Internal("internal", errors.New("failed")),
			out:  http.StatusInternalServerError,
		},
		{
			name: "400 error",
			in:   errorhandling.InvalidRequest("invalid_request", errors.New("failed")),
			out:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			encodeError(kitlog.NewNopLogger())(context.Background(), tt.in, w)

			if !assert.Equal(t, tt.out, w.Code) {
				t.Errorf("encodeError() expected HTTP StatusCode = %d, got = %d", tt.out, w.Code)
			}
		})
	}
}
