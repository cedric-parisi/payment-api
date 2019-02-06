package payments

import (
	"context"

	"github.com/cedric-parisi/payment-api/internal/models"
	"github.com/cedric-parisi/payment-api/pkg/utils"
)

// PaymentRepository create/read/update or delete on the storage
type PaymentRepository interface {
	InsertPayment(ctx context.Context, payment *models.Payment) error
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	GetFilteredPayments(ctx context.Context, filter *utils.Filter) ([]*models.Payment, int, error)
	DeletePayment(ctx context.Context, id string) error
}
