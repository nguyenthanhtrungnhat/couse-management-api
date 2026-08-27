package payment

type CreatePaymentRequest struct {
	CourseID string `json:"course_id" validate:"required,uuid"`
}