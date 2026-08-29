package seeders

import "log"

func SeedAll() {
	log.Println("🌱 Starting database seeding...")

	if err := CleanDatabase(); err != nil {
		log.Printf("❌ Database cleanup failed: %v", err)
		return
	}

	// 1. Base data
	SeedRoles()
	SeedCategories()

	// 2. Users
	SeedUsers()

	// 3. Course structure
	SeedCourses()
	SeedCourseSections()
	SeedLessons()
	SeedFileMaterials()

	// 4. Learning data
	SeedEnrollments()
	SeedProgress()

	// 5. Social data
	SeedReviews()
	SeedComments()

	// 6. Payment data
	SeedPayments()
	SeedPaymentLogs()

	log.Println("🎉 Database seeding completed successfully!")

}
