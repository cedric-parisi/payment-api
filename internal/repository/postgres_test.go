package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/cedric-parisi/payment-api/internal/models"
	"github.com/jinzhu/gorm"
)

func Test_paymentRepository_InsertPayment(t *testing.T) {
	type args struct {
		ctx     context.Context
		payment *models.Payment
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		mockCalls func(m sqlmock.Sqlmock)
	}{
		// TODO
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, _ := sqlmock.New()
			db, _ := gorm.Open("postgres", mockDB)
			p := NewPaymentRepository(db)

			tt.mockCalls(mock)

			if err := p.InsertPayment(tt.args.ctx, tt.args.payment); (err != nil) != tt.wantErr {
				t.Errorf("paymentRepository.InsertPayment() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
