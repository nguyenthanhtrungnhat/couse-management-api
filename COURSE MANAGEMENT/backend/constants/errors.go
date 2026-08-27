package constants

import "errors"

var (

	// Common
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("permission denied")
	ErrNotFound       = errors.New("resource not found")

	// Auth
	ErrEmailExists       = errors.New("email already exists")
	ErrInvalidCredential = errors.New("invalid email or password")
	ErrGoogleAccount     = errors.New("this account uses google login")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrMissingToken      = errors.New("authorization header is required")

	// User
	ErrUserNotFound = errors.New("user not found")

	// Course
	ErrCourseNotFound = errors.New("course not found")
	ErrCourseRejected = errors.New("course has been rejected")
	ErrCourseDraft    = errors.New("course is still in draft")

	// Enrollment
	ErrAlreadyEnrolled = errors.New("user already enrolled")
	ErrNotEnrolled     = errors.New("user is not enrolled")

	// Review
	ErrAlreadyReviewed = errors.New("course already reviewed")

	// Payment
	ErrPaymentFailed  = errors.New("payment failed")
	ErrPaymentPending = errors.New("payment pending")
)