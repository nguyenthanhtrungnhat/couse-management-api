
package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedReviews() {
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
		`SELECT id
		 FROM courses
		 WHERE status = $1
		 ORDER BY created_at ASC`,
		"published",
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
		log.Println("❌ No published courses found")
		return
	}

	reviews := []struct {
		studentIndex int
		courseIndex  int
		rating       int
		comment      string
	}{
		{
			studentIndex: 0,
			courseIndex:  0,
			rating:       5,
			comment:      "Excellent course. The explanations are clear and practical.",
		},
		{
			studentIndex: 0,
			courseIndex:  1,
			rating:       4,
			comment:      "Good introduction to PostgreSQL.",
		},
	}

	for _, review := range reviews {
		studentID := studentIDs[review.studentIndex%len(studentIDs)]
		courseID := courseIDs[review.courseIndex%len(courseIDs)]

		_, err := config.DB.Exec(
			ctx,
			`INSERT INTO reviews (
				id,
				user_id,
				course_id,
				rating,
				comment,
				created_at,
				updated_at
			)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
			uuid.New(),
			studentID,
			courseID,
			review.rating,
			review.comment,
		)

		if err != nil {
			log.Printf("❌ Review: %v", err)
			continue
		}

		log.Printf(
			"✅ Review: student %s -> course %s",
			studentID,
			courseID,
		)
	}
}

