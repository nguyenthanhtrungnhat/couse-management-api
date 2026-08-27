package seeders

import (
	"course-management/config"
	"course-management/models"
	"log"
)

func SeedPaymentLogs() {
	var payment models.Payment

	if err := config.DB.
		Where("transaction_code = ?", "DEMO-TXN-0001").
		First(&payment).Error; err != nil {
		return
	}

	var existing models.PaymentLog

	if err := config.DB.
		Where("payment_id = ?", payment.ID).
		First(&existing).Error; err == nil {
		return
	}

	logEntry := models.PaymentLog{
		PaymentID: payment.ID,
		RawResponse: `{
			"transaction_code": "DEMO-TXN-0001",
			"status": "success",
			"message": "Demo payment"
		}`,
	}

	if err := config.DB.Create(&logEntry).Error; err != nil {
		log.Printf("❌ Payment log: %v", err)
		return
	}

	log.Println("✅ Payment log created")
}
