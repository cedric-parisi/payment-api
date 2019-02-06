// +build !integration

package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestPayment_Validate(t *testing.T) {
	pID := uuid.New()
	tests := []struct {
		name    string
		payment *Payment
		wantErr bool
	}{
		{
			name: "payment is valid",
			payment: &Payment{
				ID:   pID,
				Type: PaymentType,
				Attribute: &Attribute{
					PaymentID: pID,
				},
			},
		},
		{
			name: "payment is invalid due to wrong type",
			payment: &Payment{
				ID:   pID,
				Type: "wrong_type",
			},
			wantErr: true,
		},
		{
			name: "payment is invalid due to missing attribute",
			payment: &Payment{
				ID:   pID,
				Type: PaymentType,
			},
			wantErr: true,
		},
		{
			name: "payment is invalid due to unrelated attribute",
			payment: &Payment{
				ID:   pID,
				Type: PaymentType,
				Attribute: &Attribute{
					PaymentID: uuid.New(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.payment.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Payment.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
