
package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedPayments() {
	ctx := context.Background()

	// Get students.
	rows, err := config.DB.Query(
		ctx,
		`SELECT u.id
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE r.name = $1
		 ORDER BY u.created_at ASC`,
		"student",
	)
	if err != nil {
		log.Printf("❌ Cannot load students: %v", err)
		return
	}

	var studentIDs []uuid.UUID

	for rows.Next() {
		var id uuid.UUID

		if err := rows.Scan(&id); err != nil {
			log.Printf("❌ Cannot read student ID: %v", err)
			continue
		}

		studentIDs = append(studentIDs, id)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		log.Printf("❌ Student iteration error: %v", err)
		return
	}

	if len(studentIDs) == 0 {
		log.Println("❌ No students found")
		return
	}

	// Get published courses.
	courseRows, err := config.DB.Query(
		ctx,
		`SELECT id, price
		 FROM courses
		 WHERE status = $1
		 ORDER BY created_at ASC`,
		"published",
	)
	if err != nil {
		log.Printf("❌ Cannot load published courses: %v", err)
		return
	}

	type courseData struct {
		id    uuid.UUID
		price float64
	}

	var courses []courseData

	for courseRows.Next() {
		var c courseData

		if err := courseRows.Scan(&c.id, &c.price); err != nil {
			log.Printf("❌ Cannot read course: %v", err)
			continue
		}

		courses = append(courses, c)
	}

	courseRows.Close()

	if err := courseRows.Err(); err != nil {
		log.Printf("❌ Course iteration error: %v", err)
		return
	}

	if len(courses) == 0 {
		log.Println("❌ No published courses found")
		return
	}

	studentID := studentIDs[0]
	c := courses[0]

	paymentID := uuid.New()

	_, err = config.DB.Exec(
		ctx,
		`INSERT INTO payments (
			id,
			user_id,
			course_id,
			amount,
			bank_name,
			transaction_code,
			status,
			paid_at,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), NOW()
		)`,
		paymentID,
		studentID,
		c.id,
		c.price,
		"Vietcombank",
		"DEMO-TXN-001",
		"success",
	)

	if err != nil {
		log.Printf("❌ Payment: %v", err)
		return
	}

	log.Printf("✅ Payment: %s", paymentID)
}

