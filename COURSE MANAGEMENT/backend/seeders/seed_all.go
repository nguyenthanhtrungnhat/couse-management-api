package seeders

func SeedAll() {
    // 0. Clean old seed data
    CleanDatabase()

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
}
