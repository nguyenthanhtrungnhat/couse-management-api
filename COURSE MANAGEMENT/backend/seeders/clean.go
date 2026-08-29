package seeders

import (
	"context"
	"fmt"
	"log"

	"course-management/config"
)

func CleanDatabase() {
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
			log.Printf("❌ Failed to clean %s: %v", table, err)
			return
		}

		log.Printf("🧹 Cleaned: %s", table)
	}

	log.Println("✅ Database cleaned successfully")
}

