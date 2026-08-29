package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedPayments() {
	log.Println("🌱 Seeding payments...")

	// Get student
	var studentID uuid.UUID

	err := config.DB.QueryRow(
		context.Background(),
		`
		SELECT id
		FROM users
		WHERE email = $1
		LIMIT 1
		`,
		"student@example.com",
	).Scan(&studentID)

	if err != nil {
		log.Printf("❌ Payment: failed to find student: %v", err)
		return
	}

	// Get first course
	var courseID uuid.UUID
	var coursePrice float64

	err = config.DB.QueryRow(
		context.Background(),
		`
		SELECT id, price
		FROM courses
		ORDER BY created_at
		LIMIT 1
		`,
	).Scan(
		&courseID,
		&coursePrice,
	)

	if err != nil {
		log.Printf("❌ Payment: failed to find course: %v", err)
		return
	}

	paymentID := uuid.New()
	transactionCode := "DEMO-TXN-" + uuid.New().String()

	_, err = config.DB.Exec(
		context.Background(),
		`
		INSERT INTO payments (
			id,
			created_at,
			updated_at,
			user_id,
			course_id,
			amount,
			bank_name,
			transaction_code,
			status
		)
		VALUES (
			$1,
			NOW(),
			NOW(),
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)
		`,
		paymentID,
		studentID,
		courseID,
		coursePrice,
		"Demo Bank",
		transactionCode,
		"success",
	)

	if err != nil {
		log.Printf("❌ Payment: %v", err)
		return
	}

	log.Printf(
		"✅ Payment: student %s -> course %s | amount %.2f | transaction %s",
		studentID,
		courseID,
		coursePrice,
		transactionCode,
	)

	log.Println("✅ Payment seeding completed")
}
