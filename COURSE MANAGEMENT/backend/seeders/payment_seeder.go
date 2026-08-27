package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedPayments() {
	var student models.User
	var course models.Course

	if err := config.DB.
		Where("email = ?", "student@example.com").
		First(&student).Error; err != nil {
		return
	}

	if err := config.DB.
		Order("created_at ASC").
		First(&course).Error; err != nil {
		return
	}

	var existing models.Payment

	if err := config.DB.
		Where("transaction_code = ?", "DEMO-TXN-0001").
		First(&existing).Error; err == nil {
		return
	}

	payment := models.Payment{
		UserID:          student.ID,
		CourseID:        course.ID,
		Amount:          course.Price,
		BankName:        "Demo Bank",
		TransactionCode: "DEMO-TXN-0001",
		Status:          "success",
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		log.Printf("❌ Payment: %v", err)
		return
	}

	log.Println("✅ Payment created")
}
