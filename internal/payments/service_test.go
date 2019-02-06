// +build !integration

package payments

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/cedric-parisi/payment-api/pkg/utils"

	"github.com/cedric-parisi/payment-api/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_service_CreatePayment(t *testing.T) {
	type args struct {
		ctx     context.Context
		payment *models.Payment
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		mockCalls func(m *MockPaymentRepository)
	}{
		{
			name: "create payment success",
			args: args{
				ctx: context.Background(),
				payment: &models.Payment{
					Type: models.PaymentType,
					Attribute: &models.Attribute{
						BeneficiaryParty: &models.BeneficiaryParty{},
						ChargesInformation: &models.ChargesInformation{
							SenderCharges: []*models.SenderCharge{{Currency: "GBP"}},
						},
						DebtorParty:  &models.DebtorParty{},
						Fx:           &models.Fx{},
						SponsorParty: &models.SponsorParty{},
					},
				},
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("InsertPayment", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "create payment failed due to invalid payment",
			args: args{
				ctx: context.Background(),
				payment: &models.Payment{
					Type: "unknown payment type",
				},
			},
			mockCalls: func(m *MockPaymentRepository) {},
			wantErr:   true,
		},
		{
			name: "create payment failed due to repository error",
			args: args{
				ctx: context.Background(),
				payment: &models.Payment{
					Type:      models.PaymentType,
					Attribute: &models.Attribute{},
				},
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("InsertPayment", mock.Anything, mock.Anything).Return(errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockPaymentRepository{}
			tt.mockCalls(mockRepo)
			s := &service{
				repository: mockRepo,
			}

			// Act
			got, err := s.CreatePayment(tt.args.ctx, tt.args.payment)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("service.CreatePayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.CreatedAt)
				assert.Equal(t, got.ID, got.Attribute.PaymentID)
			}
			assert.True(t, mock.AssertExpectationsForObjects(t, mockRepo))
		})
	}
}

func Test_service_UpdatePayment(t *testing.T) {
	pID := uuid.New()
	type args struct {
		ctx     context.Context
		payment *models.Payment
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		mockCalls func(m *MockPaymentRepository)
	}{
		{
			name: "update payment success",
			args: args{
				ctx: context.Background(),
				payment: &models.Payment{
					ID:   pID,
					Type: models.PaymentType,
					Attribute: &models.Attribute{
						PaymentID: pID,
					},
				},
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("UpdatePayment", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "update payment failed due to invalid payment",
			args: args{
				ctx: context.Background(),
				payment: &models.Payment{
					ID:   pID,
					Type: "unknown payment type",
				},
			},
			mockCalls: func(m *MockPaymentRepository) {},
			wantErr:   true,
		},
		{
			name: "update payment failed due to repository error",
			args: args{
				ctx: context.Background(),
				payment: &models.Payment{
					ID:   pID,
					Type: models.PaymentType,
					Attribute: &models.Attribute{
						PaymentID: pID,
					},
				},
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("UpdatePayment", mock.Anything, mock.Anything).Return(errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockPaymentRepository{}
			tt.mockCalls(mockRepo)
			s := &service{
				repository: mockRepo,
			}
			// Act
			err := s.UpdatePayment(tt.args.ctx, tt.args.payment)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("service.UpdatePayment() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.True(t, mock.AssertExpectationsForObjects(t, mockRepo))
		})
	}
}

func Test_service_GetPayment(t *testing.T) {
	pID := uuid.New()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name      string
		args      args
		want      *models.Payment
		wantErr   bool
		mockCalls func(m *MockPaymentRepository)
	}{
		{
			name: "get payment success",
			args: args{
				ctx: context.Background(),
				id:  pID.String(),
			},
			want: &models.Payment{
				ID: pID,
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("GetPayment", mock.Anything, mock.Anything).Return(&models.Payment{
					ID: pID,
				}, nil)
			},
		},
		{
			name: "get payment failed due to invalid id",
			args: args{
				ctx: context.Background(),
				id:  "invalid id",
			},
			mockCalls: func(m *MockPaymentRepository) {},
			wantErr:   true,
		},
		{
			name: "get payment not found",
			args: args{
				ctx: context.Background(),
				id:  pID.String(),
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("GetPayment", mock.Anything, mock.Anything).Return(nil, ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "get payment failed due to repository error",
			args: args{
				ctx: context.Background(),
				id:  pID.String(),
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("GetPayment", mock.Anything, mock.Anything).Return(nil, errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockPaymentRepository{}
			tt.mockCalls(mockRepo)
			s := &service{
				repository: mockRepo,
			}

			// Act
			got, err := s.GetPayment(tt.args.ctx, tt.args.id)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetPayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetPayment() = %v, want %v", got, tt.want)
			}
			assert.True(t, mock.AssertExpectationsForObjects(t, mockRepo))
		})
	}
}

func Test_service_GetFilteredPayments(t *testing.T) {
	type args struct {
		ctx     context.Context
		filters *utils.Filter
	}
	tests := []struct {
		name      string
		args      args
		want      *utils.FilteredList
		wantErr   bool
		mockCalls func(m *MockPaymentRepository)
	}{
		{
			name: "get filtered payments success",
			args: args{
				ctx: context.Background(),
				filters: &utils.Filter{
					Limit:  10,
					Offset: 0,
				},
			},
			want: &utils.FilteredList{
				Resource: resourceName,
				Results:  []*models.Payment{},
				Filter: utils.Filter{
					Limit:  10,
					Offset: 0,
				},
				TotalCount: 10,
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("GetFilteredPayments", mock.Anything, mock.Anything).Return([]*models.Payment{}, 10, nil)
			},
		},
		{
			name: "get filtered failed due to wrong limit",
			args: args{
				ctx: context.Background(),
				filters: &utils.Filter{
					Limit: maxLimit + 10,
				},
			},
			mockCalls: func(m *MockPaymentRepository) {},
			wantErr:   true,
		},
		{
			name: "get filtered failed due to repository error",
			args: args{
				ctx: context.Background(),
				filters: &utils.Filter{
					Limit:  10,
					Offset: 0,
				},
			},
			mockCalls: func(m *MockPaymentRepository) {
				m.On("GetFilteredPayments", mock.Anything, mock.Anything).Return(nil, 0, errors.New("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockPaymentRepository{}
			tt.mockCalls(mockRepo)
			s := &service{
				repository: mockRepo,
			}

			// Act
			got, err := s.GetFilteredPayments(tt.args.ctx, tt.args.filters)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetFilteredPayments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetFilteredPayments() = %v, want %v", got, tt.want)
			}
			assert.True(t, mock.AssertExpectationsForObjects(t, mockRepo))
		})
	}
}
