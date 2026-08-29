package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedPaymentLogs() {
	ctx := context.Background()

	rows, err := config.DB.Query(
		ctx,
		`SELECT id, transaction_code, status
		 FROM payments
		 ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load payments: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			paymentID       uuid.UUID
			transactionCode string
			status          string
		)

		if err := rows.Scan(
			&paymentID,
			&transactionCode,
			&status,
		); err != nil {
			log.Printf("❌ Cannot read payment: %v", err)
			continue
		}

		rawResponse := []byte(`{
			"transaction_code": "` + transactionCode + `",
			"status": "` + status + `",
			"source": "seed"
		}`)

		_, err := config.DB.Exec(
			ctx,
			`INSERT INTO payment_logs (
				id,
				payment_id,
				raw_response,
				created_at
			)
			VALUES ($1, $2, $3::jsonb, NOW())`,
			uuid.New(),
			paymentID,
			rawResponse,
		)

		if err != nil {
			log.Printf("❌ Payment log: %v", err)
			continue
		}

		log.Printf("✅ Payment log for payment: %s", paymentID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Payment iteration error: %v", err)
	}
}

