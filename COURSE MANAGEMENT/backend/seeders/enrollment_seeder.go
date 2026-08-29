
package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedEnrollments() {
	ctx := context.Background()

	userRows, err := config.DB.Query(
		ctx,
		`SELECT id FROM users ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load users: %v", err)
		return
	}

	var userIDs []uuid.UUID

	for userRows.Next() {
		var id uuid.UUID

		if err := userRows.Scan(&id); err != nil {
			log.Printf("❌ Cannot read user ID: %v", err)
			continue
		}

		userIDs = append(userIDs, id)
	}

	userRows.Close()

	if err := userRows.Err(); err != nil {
		log.Printf("❌ User iteration error: %v", err)
		return
	}

	if len(userIDs) == 0 {
		log.Println("❌ No users found")
		return
	}

	courseRows, err := config.DB.Query(
		ctx,
		`SELECT id FROM courses ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load courses: %v", err)
		return
	}

	var courseIDs []uuid.UUID

	for courseRows.Next() {
		var id uuid.UUID

		if err := courseRows.Scan(&id); err != nil {
			log.Printf("❌ Cannot read course ID: %v", err)
			continue
		}

		courseIDs = append(courseIDs, id)
	}

	courseRows.Close()

	if err := courseRows.Err(); err != nil {
		log.Printf("❌ Course iteration error: %v", err)
		return
	}

	if len(courseIDs) == 0 {
		log.Println("❌ No courses found")
		return
	}

	// Find a student instead of using the instructor account.
	var studentID uuid.UUID

	err = config.DB.QueryRow(
		ctx,
		`SELECT u.id
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE r.name = $1
		 ORDER BY u.created_at ASC
		 LIMIT 1`,
		"student",
	).Scan(&studentID)

	if err != nil {
		log.Printf("❌ Cannot find student: %v", err)
		return
	}

	// Enroll the student in the first published course.
	var courseID uuid.UUID

	err = config.DB.QueryRow(
		ctx,
		`SELECT id
		 FROM courses
		 WHERE status = $1
		 ORDER BY created_at ASC
		 LIMIT 1`,
		"published",
	).Scan(&courseID)

	if err != nil {
		log.Printf("❌ Cannot find published course: %v", err)
		return
	}

	_, err = config.DB.Exec(
		ctx,
		`INSERT INTO enrollments (
			id,
			user_id,
			course_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, NOW(), NOW(), NOW())`,
		uuid.New(),
		studentID,
		courseID,
	)

	if err != nil {
		log.Printf("❌ Enrollment: %v", err)
		return
	}

	log.Printf(
		"✅ Enrollment: student %s -> course %s",
		studentID,
		courseID,
	)
}

