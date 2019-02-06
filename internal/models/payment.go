package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Type is payment type
type Type string

const (
	// PaymentType ...
	PaymentType = "Payment"
	// WithdrawType ...
	WithdrawType = "Withdraw"
)

// Payment define a payment
type Payment struct {
	ID             uuid.UUID  `json:"id" gorm:"primary_key"`
	Type           Type       `json:"type"`
	Version        int        `json:"version"`
	OrganisationID uuid.UUID  `json:"organisation_id"`
	Attribute      *Attribute `json:"attributes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"-"`
}

// Attribute ...
type Attribute struct {
	ID                   uuid.UUID           `json:"id" gorm:"primary_key"`
	PaymentID            uuid.UUID           `json:"payment_id" gorm:"foreign_key"`
	Amount               string              `json:"amount"`
	BeneficiaryParty     *BeneficiaryParty   `json:"beneficiary_party"`
	ChargesInformation   *ChargesInformation `json:"charges_information"`
	Currency             string              `json:"currency"`
	DebtorParty          *DebtorParty        `json:"debtor_party"`
	EndToEndReference    string              `json:"end_to_end_reference"`
	Fx                   *Fx                 `json:"fx"`
	NumericReference     string              `json:"numeric_reference"`
	PaymentPurpose       string              `json:"payment_purpose"`
	PaymentScheme        string              `json:"payment_scheme"`
	PaymentType          string              `json:"payment_type"`
	ProcessingDate       string              `json:"processing_date"`
	Reference            string              `json:"reference"`
	SchemePaymentSubType string              `json:"scheme_payment_sub_type"`
	SchemePaymentType    string              `json:"scheme_payment_type"`
	SponsorParty         *SponsorParty       `json:"sponsor_party"`
}

// BeneficiaryParty ...
type BeneficiaryParty struct {
	AttributeID       uuid.UUID `json:"attribute_id" gorm:"foreign_key"`
	AccountName       string    `json:"account_name"`
	AccountNumber     string    `json:"account_number"`
	AccountNumberCode string    `json:"account_number_code"`
	AccountType       int       `json:"account_type"`
	Address           string    `json:"address"`
	BankID            string    `json:"bank_id"`
	BankIDCode        string    `json:"bank_id_code"`
	Name              string    `json:"name"`
}

// ChargesInformation ...
type ChargesInformation struct {
	ID                      uuid.UUID       `json:"id" gorm:"primary_key"`
	AttributeID             uuid.UUID       `json:"attribute_id" gorm:"foreign_key"`
	BearerCode              string          `json:"bearer_code"`
	SenderCharges           []*SenderCharge `json:"sender_charges"`
	ReceiverChargesAmount   string          `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `json:"receiver_charges_currency"`
}

// SenderCharge ...
type SenderCharge struct {
	ChargesInformationID uuid.UUID `json:"charges_information_id" gorm:"foreign_key"`
	Amount               string    `json:"amount"`
	Currency             string    `json:"currency"`
}

// DebtorParty ...
type DebtorParty struct {
	AttributeID       uuid.UUID `json:"attribute_id" gorm:"foreign_key"`
	AccountName       string    `json:"account_name"`
	AccountNumber     string    `json:"account_number"`
	AccountNumberCode string    `json:"account_number_code"`
	Address           string    `json:"address"`
	BankID            string    `json:"bank_id"`
	BankIDCode        string    `json:"bank_id_code"`
	Name              string    `json:"name"`
}

// Fx ...
type Fx struct {
	AttributeID       uuid.UUID `json:"attribute_id" gorm:"foreign_key"`
	ContractReference string    `json:"contract_reference"`
	ExchangeRate      string    `json:"exchange_rate"`
	OriginalAmount    string    `json:"original_amount"`
	OriginalCurrency  string    `json:"original_currency"`
}

// SponsorParty ...
type SponsorParty struct {
	AttributeID   uuid.UUID `json:"attribute_id" gorm:"foreign_key"`
	AccountNumber string    `json:"account_number"`
	BankID        string    `json:"bank_id"`
	BankIDCode    string    `json:"bank_id_code"`
}

// Validate ensures that the payment is valid
func (p *Payment) Validate() error {
	if p.Type != PaymentType && p.Type != WithdrawType {
		return errors.New("invalid payment type")
	}

	if p.Attribute == nil {
		return errors.New("attribute is required")
	}

	if p.ID != p.Attribute.PaymentID {
		return errors.New("attribute and payment not related")
	}

	// TODO validate all fields

	return nil
}
