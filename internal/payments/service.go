package payments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/cedric-parisi/payment-api/pkg/utils"

	"github.com/cedric-parisi/payment-api/internal/models"
	"github.com/cedric-parisi/payment-api/pkg/errorhandling"
)

const (
	resourceName = "payments"
	maxLimit     = 500

	invalidPaymentCode    = "invalid_payment"
	persistFailedCode     = "save_payment_failed"
	readPaymentFailedCode = "read_payment_failed"
)

var (
	// ErrNotFound is raised when a payment is not found in the storage
	ErrNotFound = errors.New("not found")
)

// Service defines the business logic on the payment resource
type Service interface {
	CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	GetFilteredPayments(ctx context.Context, filter *utils.Filter) (*utils.FilteredList, error)
	DeletePayment(ctx context.Context, id string) error
}

type service struct {
	repository PaymentRepository
}

// NewService ...
func NewService(repo PaymentRepository) Service {
	return &service{
		repository: repo,
	}
}

// CreatePayment creates a new payment
// Returns the newly created payment
func (s *service) CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	// Complete payment request with ids and date
	paymentID := uuid.New()
	attributeID := uuid.New()
	chargeInformationID := uuid.New()

	payment.ID = paymentID
	payment.CreatedAt = time.Now().UTC()
	if payment.Attribute != nil {
		payment.Attribute.ID = attributeID
		payment.Attribute.PaymentID = paymentID

		if payment.Attribute.BeneficiaryParty != nil {
			payment.Attribute.BeneficiaryParty.AttributeID = attributeID
		}

		if payment.Attribute.ChargesInformation != nil {
			payment.Attribute.ChargesInformation.ID = chargeInformationID
			payment.Attribute.ChargesInformation.AttributeID = attributeID
			for _, s := range payment.Attribute.ChargesInformation.SenderCharges {
				s.ChargesInformationID = chargeInformationID
			}
		}

		if payment.Attribute.DebtorParty != nil {
			payment.Attribute.DebtorParty.AttributeID = attributeID
		}

		if payment.Attribute.Fx != nil {
			payment.Attribute.Fx.AttributeID = attributeID
		}

		if payment.Attribute.SponsorParty != nil {
			payment.Attribute.SponsorParty.AttributeID = attributeID
		}
	}

	if err := payment.Validate(); err != nil {
		return nil, errorhandling.InvalidRequest(invalidPaymentCode, err)
	}

	if err := s.repository.InsertPayment(ctx, payment); err != nil {
		return nil, errorhandling.Internal(persistFailedCode, err)
	}
	return payment, nil
}

// UpdatePayment updates an existing payment
func (s *service) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	if err := payment.Validate(); err != nil {
		return errorhandling.InvalidRequest(invalidPaymentCode, err)
	}

	if err := s.repository.UpdatePayment(ctx, payment); err != nil {
		return errorhandling.Internal(persistFailedCode, err)
	}
	return nil
}

// GetPayment returns the payment resource selected by its unique identifier
func (s *service) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errorhandling.InvalidRequest(invalidPaymentCode, err)
	}

	payment, err := s.repository.GetPayment(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			return nil, errorhandling.NotFound(invalidPaymentCode, fmt.Errorf("could not find %s", id))
		}
		return nil, errorhandling.Internal(readPaymentFailedCode, err)
	}

	return payment, nil
}

// GetFilteredPayments returns a list of payments matching the requested filters
func (s *service) GetFilteredPayments(ctx context.Context, filters *utils.Filter) (*utils.FilteredList, error) {
	if filters.Limit > maxLimit {
		return nil, errorhandling.InvalidRequest(invalidPaymentCode, fmt.Errorf("limit must be lower than %d", maxLimit))
	}

	payments, totalCount, err := s.repository.GetFilteredPayments(ctx, filters)
	if err != nil {
		return nil, errorhandling.Internal(readPaymentFailedCode, err)
	}
	return &utils.FilteredList{
		Results:    payments,
		Resource:   resourceName,
		Filter:     *filters,
		TotalCount: totalCount,
	}, nil
}

// DeletePayment deletes the payment by its unique identifier
func (s *service) DeletePayment(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errorhandling.InvalidRequest(invalidPaymentCode, err)
	}

	if err := s.repository.DeletePayment(ctx, id); err != nil {
		return errorhandling.Internal(persistFailedCode, err)
	}
	return nil
}
