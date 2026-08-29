package seeders

import (
"context"
"log"


"course-management/config"
"course-management/utils"

"github.com/google/uuid"


)

func SeedUsers() {
ctx := context.Background()

users := []struct {
	Name     string
	Email    string
	Role     string
	Password string
}{
	{
		Name:     "System Admin",
		Email:    "admin@example.com",
		Role:     "admin",
		Password: "123456",
	},
	{
		Name:     "Go Instructor",
		Email:    "instructor@example.com",
		Role:     "instructor",
		Password: "123456",
	},
	{
		Name:     "Web Instructor",
		Email:    "webinstructor@example.com",
		Role:     "instructor",
		Password: "123456",
	},
	{
		Name:     "Demo Student",
		Email:    "student@example.com",
		Role:     "student",
		Password: "123456",
	},
	{
		Name:     "Test Student",
		Email:    "student2@example.com",
		Role:     "student",
		Password: "123456",
	},
	{
		Name:     "Another Student",
		Email:    "student3@example.com",
		Role:     "student",
		Password: "123456",
	},
}

for _, data := range users {
	var roleID uuid.UUID

	err := config.DB.QueryRow(
		ctx,
		`SELECT id FROM roles WHERE name = $1`,
		data.Role,
	).Scan(&roleID)

	if err != nil {
		log.Printf("❌ Role %s not found: %v", data.Role, err)
		continue
	}

	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		log.Printf("❌ Password hash failed for %s: %v", data.Email, err)
		continue
	}

	_, err = config.DB.Exec(
		ctx,
		`INSERT INTO users
			(id, role_id, full_name, email, password_hash, provider, created_at, updated_at)
		 VALUES
			($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		uuid.New(),
		roleID,
		data.Name,
		data.Email,
		hash,
		"local",
	)

	if err != nil {
		log.Printf("❌ Failed user %s: %v", data.Email, err)
		continue
	}

	log.Printf("✅ User: %s (%s)", data.Email, data.Role)
}


}
