
package models

import "github.com/google/uuid"

type Payment struct {
	BaseModel

	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`

	Amount          float64 `json:"amount"`
	BankName        *string `json:"bank_name,omitempty"`
	TransactionCode *string `json:"transaction_code,omitempty"`
	Status          string  `json:"status"`

	User   *User   `json:"user,omitempty"`
	Course *Course `json:"course,omitempty"`

	PaymentLogs []PaymentLog `json:"payment_logs,omitempty"`
}

// TableName returns the database table name.
func (Payment) TableName() string {
	return "payments"
}

// NewPayment creates a new payment.
func NewPayment(
	userID uuid.UUID,
	courseID uuid.UUID,
	amount float64,
	bankName *string,
	transactionCode *string,
	status string,
) *Payment {
	return &Payment{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},

		UserID:          userID,
		CourseID:        courseID,
		Amount:          amount,
		BankName:        bankName,
		TransactionCode: transactionCode,
		Status:          status,
	}
}

