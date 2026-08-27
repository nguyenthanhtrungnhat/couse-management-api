package payment

import "time"

type PaymentResponse struct {
	ID              string    `json:"id"`
	CourseID        string    `json:"course_id"`
	Amount          float64   `json:"amount"`
	BankName        string    `json:"bank_name"`
	TransactionCode string    `json:"transaction_code"`
	Status          string    `json:"status"`
	PaidAt          time.Time `json:"paid_at"`
}