package models

import "github.com/google/uuid"

type PaymentLog struct {
	BaseModel

	PaymentID uuid.UUID `gorm:"type:uuid;index;not null"`

	Payment Payment

	RawResponse string `gorm:"type:jsonb"`
}