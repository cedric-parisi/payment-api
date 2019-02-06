// +build !integration

package payments

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/cedric-parisi/payment-api/pkg/utils"

	"github.com/cedric-parisi/payment-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestCreatePaymentResponse_StatusCode(t *testing.T) {
	// Arrange
	c := CreatePaymentResponse{}

	// Act
	got := c.StatusCode()

	// Assert
	assert.Equal(t, http.StatusCreated, got)
}

func TestCreatePaymentResponse_Headers(t *testing.T) {
	// Arrange
	cID := uuid.New()
	expectedHeader := http.Header{
		"Location": []string{fmt.Sprintf("/payments/%s", cID)},
	}
	c := CreatePaymentResponse{
		Payment: models.Payment{
			ID: cID,
		},
	}

	// Act
	got := c.Headers()

	// Assert
	assert.Equal(t, expectedHeader, got)
}

func TestMakeCreatePaymentEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		mockCalls func(svc *MockService)
		want      CreatePaymentResponse
		wantErr   bool
	}{
		{
			name: "create endpoint ok",
			mockCalls: func(svc *MockService) {
				svc.On("CreatePayment", mock.Anything, mock.Anything).Return(&models.Payment{}, nil)
			},
			want:    CreatePaymentResponse{},
			wantErr: false,
		},
		{
			name: "create endpoint returned error",
			mockCalls: func(svc *MockService) {
				svc.On("CreatePayment", mock.Anything, mock.Anything).Return(nil, errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &MockService{}
			tt.mockCalls(svc)
			endpoint := MakeCreatePaymentEndpoint(svc)
			resp, err := endpoint(context.Background(), &models.Payment{})

			// Assert
			assert.NotNil(t, endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeCreatePaymentEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, resp)
			}
		})
	}

}

func TestMakeUpdatePaymentEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		mockCalls func(svc *MockService)
		wantErr   bool
	}{
		{
			name: "update endpoint ok",
			mockCalls: func(svc *MockService) {
				svc.On("UpdatePayment", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "update endpoint returned error",
			mockCalls: func(svc *MockService) {
				svc.On("UpdatePayment", mock.Anything, mock.Anything).Return(errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &MockService{}
			tt.mockCalls(svc)
			endpoint := MakeUpdatePaymentEndpoint(svc)
			resp, err := endpoint(context.Background(), &models.Payment{})

			// Assert
			assert.NotNil(t, endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeUpdatePaymentEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, resp)
			}
		})
	}
}

func TestMakeGetPaymentEndpoint(t *testing.T) {
	svc := &MockService{}
	svc.On("GetPayment", mock.Anything, mock.Anything).Return(&models.Payment{}, nil)

	endpoint := MakeGetPaymentEndpoint(svc)
	a, err := endpoint(context.Background(), "id")

	assert.NotNil(t, endpoint)
	assert.Nil(t, err)
	assert.NotNil(t, a)

}

func TestMakeGetFilteredPaymentsEndpoint(t *testing.T) {
	svc := &MockService{}
	svc.On("GetFilteredPayments", mock.Anything, mock.Anything).Return(&utils.FilteredList{}, nil)

	endpoint := MakeGetFilteredPaymentsEndpoint(svc)
	a, err := endpoint(context.Background(), &utils.Filter{})

	assert.NotNil(t, endpoint)
	assert.Nil(t, err)
	assert.NotNil(t, a)

}

func TestMakeDeletePaymentEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		mockCalls func(svc *MockService)
		wantErr   bool
	}{
		{
			name: "delete endpoint ok",
			mockCalls: func(svc *MockService) {
				svc.On("DeletePayment", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete endpoint returned error",
			mockCalls: func(svc *MockService) {
				svc.On("DeletePayment", mock.Anything, mock.Anything).Return(errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &MockService{}
			tt.mockCalls(svc)
			endpoint := MakeDeletePaymentEndpoint(svc)
			resp, err := endpoint(context.Background(), "id")

			// Assert
			assert.NotNil(t, endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeDeletePaymentEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, resp)
			}
		})
	}
}
