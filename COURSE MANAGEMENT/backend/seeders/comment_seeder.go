package seeders

import (
	"context"
	"log"

	"course-management/config"

	"github.com/google/uuid"
)

func SeedComments() {
	ctx := context.Background()

	// Get users.
	rows, err := config.DB.Query(
		ctx,
		`SELECT id FROM users ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("❌ Cannot load users: %v", err)
		return
	}

	var userIDs []uuid.UUID

	for rows.Next() {
		var id uuid.UUID

		if err := rows.Scan(&id); err != nil {
			log.Printf("❌ Cannot read user ID: %v", err)
			continue
		}

		userIDs = append(userIDs, id)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		log.Printf("❌ User iteration error: %v", err)
		return
	}

	if len(userIDs) == 0 {
		log.Println("❌ No users found")
		return
	}

	// Get courses.
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

	comments := []struct {
		userIndex   int
		courseIndex int
		content     string
	}{
		{
			userIndex:   0,
			courseIndex: 0,
			content:     "This course is very easy to follow.",
		},
		{
			userIndex:   1,
			courseIndex: 0,
			content:     "The examples are useful and practical.",
		},
		{
			userIndex:   0,
			courseIndex: 1,
			content:     "Good introduction to PostgreSQL.",
		},
	}

	for _, comment := range comments {
		userID := userIDs[comment.userIndex%len(userIDs)]
		courseID := courseIDs[comment.courseIndex%len(courseIDs)]

		_, err := config.DB.Exec(
			ctx,
			`INSERT INTO comments (
				id,
				user_id,
				course_id,
				content,
				created_at,
				updated_at
			)
			VALUES ($1, $2, $3, $4, NOW(), NOW())`,
			uuid.New(),
			userID,
			courseID,
			comment.content,
		)

		if err != nil {
			log.Printf("❌ Comment: %v", err)
			continue
		}

		log.Printf("✅ Comment: %s", comment.content)
	}
}
