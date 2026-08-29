package models

import "github.com/google/uuid"

type Payment struct {
	BaseModel

	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`

	User   User   `json:"user,omitempty"`
	Course Course `json:"course,omitempty"`

	Amount          float64 `json:"amount"`
	BankName        string  `json:"bank_name"`
	TransactionCode string  `json:"transaction_code"`
	Status          string  `json:"status"`

	Logs []PaymentLog `json:"logs,omitempty"`
}