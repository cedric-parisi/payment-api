package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cedric-parisi/payment-api/internal/payments"

	"github.com/cedric-parisi/payment-api/pkg/utils"

	"github.com/jinzhu/gorm"

	"github.com/cedric-parisi/payment-api/internal/models"
)

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository ...
func NewPaymentRepository(db *gorm.DB) payments.PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

// InsertPayment save a new payment
func (p paymentRepository) InsertPayment(ctx context.Context, payment *models.Payment) error {
	return p.db.Create(payment).Error
}

// UpdatePayment updates an existing payment
func (p paymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	return p.db.Save(payment).Error
}

// GetPayment select a payment by its id
func (p paymentRepository) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	payment := &models.Payment{}
	err := p.db.First(payment, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, payments.ErrNotFound
		}
		return nil, err
	}

	err = p.getRelated(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetFilteredPayments selects payments according to filters
func (p paymentRepository) GetFilteredPayments(ctx context.Context, filter *utils.Filter) ([]*models.Payment, int, error) {
	stmt := p.db.Offset(filter.Offset).Limit(filter.Limit)
	for _, sort := range filter.Sorting {
		direction := "ASC"
		if sort.Descending {
			direction = "DESC"
		}
		stmt = stmt.Order(fmt.Sprintf("%s %s", sort.Field, direction))
	}
	var payments []*models.Payment
	totalCount := 0
	err := stmt.Find(&payments).Count(&totalCount).Error
	if err != nil {
		// Find returns sql.ErrNoRows when no result found
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	for _, payment := range payments {
		err = p.getRelated(ctx, payment)
		if err != nil {
			return nil, 0, err
		}
	}

	return payments, totalCount, nil
}

func (p paymentRepository) DeletePayment(ctx context.Context, id string) error {
	return p.db.Delete(models.Payment{}, "id = ?", id).Error
}

func (p paymentRepository) getRelated(ctx context.Context, payment *models.Payment) error {
	payment.Attribute = &models.Attribute{
		BeneficiaryParty: &models.BeneficiaryParty{},
		ChargesInformation: &models.ChargesInformation{
			SenderCharges: []*models.SenderCharge{},
		},
		DebtorParty:  &models.DebtorParty{},
		Fx:           &models.Fx{},
		SponsorParty: &models.SponsorParty{},
	}
	err := p.db.First(payment.Attribute, "payment_id = ?", payment.ID).
		Related(payment.Attribute.BeneficiaryParty).
		Related(payment.Attribute.BeneficiaryParty).
		Related(payment.Attribute.ChargesInformation).
		Related(payment.Attribute.DebtorParty).
		Related(payment.Attribute.Fx).
		Related(payment.Attribute.SponsorParty).
		Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	err = p.db.Find(&payment.Attribute.ChargesInformation.SenderCharges, "charges_information_id = ?", payment.Attribute.ChargesInformation.ID).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	return nil
}
