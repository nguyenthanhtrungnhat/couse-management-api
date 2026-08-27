package constants

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"

	PaymentSuccess PaymentStatus = "success"

	PaymentFailed PaymentStatus = "failed"
)
