package seeders

import (
"context"
"log"


"course-management/config"
"course-management/models"


)

func SeedCategories() {
ctx := context.Background()


categories := []models.Category{
	{
		Name:        "Programming",
		Description: "Programming and software development",
	},
	{
		Name:        "Web Development",
		Description: "Frontend, backend and full-stack web development",
	},
	{
		Name:        "Database",
		Description: "SQL, PostgreSQL and database management",
	},
	{
		Name:        "DevOps",
		Description: "Deployment, Docker, cloud and DevOps",
	},
	{
		Name:        "Mobile Development",
		Description: "Mobile application development",
	},
}

for _, category := range categories {
	_, err := config.DB.Exec(
		ctx,
		`INSERT INTO categories
			(id, name, description, created_at, updated_at)
		 VALUES
			(gen_random_uuid(), $1, $2, NOW(), NOW())`,
		category.Name,
		category.Description,
	)

	if err != nil {
		log.Printf("❌ Failed to seed category %s: %v", category.Name, err)
		continue
	}

	log.Printf("✅ Category: %s", category.Name)
}


}

func stringPtr(s string) *string {
return &s
}
