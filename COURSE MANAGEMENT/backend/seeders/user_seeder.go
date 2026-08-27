package seeders

import (
	"course-management/config"
	"course-management/models"
	"course-management/utils"
	"log"
)

func SeedUsers() {
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
		var existing models.User

		if err := config.DB.
			Where("email = ?", data.Email).
			First(&existing).Error; err == nil {
			continue
		}

		var role models.Role

		if err := config.DB.
			Where("name = ?", data.Role).
			First(&role).Error; err != nil {
			log.Printf("❌ Role %s not found", data.Role)
			continue
		}

		hash, err := utils.HashPassword(data.Password)
		if err != nil {
			log.Printf("❌ Password hash failed: %v", err)
			continue
		}

		user := models.User{
			RoleID:       role.ID,
			FullName:     data.Name,
			Email:        data.Email,
			PasswordHash: &hash,
			Provider:     "local",
		}

		if err := config.DB.Create(&user).Error; err != nil {
			log.Printf("❌ Failed user %s: %v", data.Email, err)
			continue
		}

		log.Printf("✅ User: %s (%s)", data.Email, data.Role)
	}
}