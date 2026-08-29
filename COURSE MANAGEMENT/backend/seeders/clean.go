package seeders

import (
	"context"
	"fmt"
	"log"

	"course-management/config"
)

func CleanDatabase() error {
	tables := []string{
		"payment_logs",
		"payments",
		"comments",
		"reviews",
		"lesson_progress",
		"enrollments",
		"file_materials",
		"lessons",
		"course_sections",
		"courses",
		"categories",
		"users",
		"roles",
	}

	ctx := context.Background()

	for _, table := range tables {
		query := fmt.Sprintf(`DELETE FROM "%s"`, table)

		if _, err := config.DB.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}

		log.Printf("🧹 Cleaned: %s", table)
	}

	log.Println("✅ Database cleaned successfully")

	return nil

}
