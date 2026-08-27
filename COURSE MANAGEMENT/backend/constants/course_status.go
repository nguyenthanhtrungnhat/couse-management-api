package constants

type CourseStatus string

const (

    CourseDraft CourseStatus = "draft"

    CoursePending CourseStatus = "pending"

    CoursePublished CourseStatus = "published"

    CourseRejected CourseStatus = "rejected"
)