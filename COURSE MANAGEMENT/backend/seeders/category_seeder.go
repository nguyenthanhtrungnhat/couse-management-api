package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedCategories() {
	categories := []models.Category{
		{
			Name:        "Programming",
			Description: stringPtr("Programming and software development"),
		},
		{
			Name:        "Web Development",
			Description: stringPtr("Frontend, backend and full-stack web development"),
		},
		{
			Name:        "Database",
			Description: stringPtr("SQL, PostgreSQL and database management"),
		},
		{
			Name:        "DevOps",
			Description: stringPtr("Deployment, Docker, cloud and DevOps"),
		},
		{
			Name:        "Mobile Development",
			Description: stringPtr("Mobile application development"),
		},
	}

	for _, category := range categories {
		var existing models.Category

		err := config.DB.
			Where("name = ?", category.Name).
			First(&existing).Error

		if err == nil {
			continue
		}

		if err := config.DB.Create(&category).Error; err != nil {
			log.Printf("❌ Failed to seed category %s: %v", category.Name, err)
			continue
		}

		log.Printf("✅ Category: %s", category.Name)
	}
}

func stringPtr(s string) *string {
	return &s
}
