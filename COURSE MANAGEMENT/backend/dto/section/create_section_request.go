package section

type CreateSectionRequest struct {
	CourseID  string `json:"course_id" validate:"required,uuid"`
	Title     string `json:"title" validate:"required,min=1,max=255"`
	SortOrder int    `json:"sort_order" validate:"min=0"`
}
