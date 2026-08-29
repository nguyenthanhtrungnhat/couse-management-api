package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type PaymentLog struct {
	BaseModel

	PaymentID   uuid.UUID       `json:"payment_id"`
	RawResponse json.RawMessage `json:"raw_response"`

	Payment *Payment `json:"payment,omitempty"`
}

// TableName returns the database table name.
func (PaymentLog) TableName() string {
	return "payment_logs"
}

// NewPaymentLog creates a new payment log.
func NewPaymentLog(
	paymentID uuid.UUID,
	rawResponse json.RawMessage,
) *PaymentLog {
	return &PaymentLog{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		PaymentID:   paymentID,
		RawResponse: rawResponse,
	}
}
