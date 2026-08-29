package models

import "github.com/google/uuid"

type PaymentLog struct {
	BaseModel

	PaymentID uuid.UUID `json:"payment_id"`
	Payment   Payment  `json:"payment,omitempty"`

	RawResponse string `json:"raw_response"`
}