package models

import "github.com/google/uuid"

type Payment struct {
	BaseModel

	UserID uuid.UUID `gorm:"type:uuid;index;not null"`

	CourseID uuid.UUID `gorm:"type:uuid;index;not null"`

	User User

	Course Course

	Amount float64 `gorm:"type:numeric(12,2)"`

	BankName string `gorm:"size:100"`

	TransactionCode string `gorm:"size:255;uniqueIndex"`

	Status string `gorm:"size:20"`

	Logs []PaymentLog `gorm:"foreignKey:PaymentID"`
}
