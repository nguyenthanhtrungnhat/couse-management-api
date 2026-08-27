package seeders

import (
	"course-management/config"
	"log"
)

func CleanDatabase() {
	log.Println("🧹 Cleaning database...")

	tables := []string{
		"payment_logs",
		"payments",

		"comments",
		"reviews",

		"progress",
		"enrollments",

		"file_materials",
		"lessons",
		"course_sections",

		"courses",

		"users",
		"categories",
		"roles",
	}

	for _, table := range tables {
		if err := config.DB.Exec(
			"DELETE FROM " + table,
		).Error; err != nil {
			log.Printf("❌ Failed to clean %s: %v", table, err)
			return
		}

		log.Printf("✅ Cleaned: %s", table)
	}

	log.Println("✅ Database cleaned successfully")
}