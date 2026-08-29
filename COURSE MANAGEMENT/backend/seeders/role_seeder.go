package seeders

import (
	"context"
	"log"

	"course-management/config"
)

func SeedRoles() {
	ctx := context.Background()

	roles := []string{
		"admin",
		"student",
		"instructor",
	}

	for _, role := range roles {
		_, err := config.DB.Exec(
			ctx,
			`INSERT INTO roles
			(id, name, created_at, updated_at)
		 VALUES
			(gen_random_uuid(), $1, NOW(), NOW())`,
			role,
		)

		if err != nil {
			log.Printf("❌ Failed to seed role %s: %v", role, err)
			continue
		}

		log.Printf("✅ Role: %s", role)
	}

}
