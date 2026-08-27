package seeders

import (
	"course-management/config"
	"course-management/models"
)

func SeedRoles() {
	var count int64

	config.DB.Model(&models.Role{}).Count(&count)

	if count > 0 {
		return
	}

	roles := []models.Role{
		{Name: "admin"},
		{Name: "student"},
		{Name: "instructor"},
	}

	config.DB.Create(&roles)
}
