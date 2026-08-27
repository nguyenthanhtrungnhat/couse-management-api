package section

type UpdateSectionRequest struct {
	Title     string `json:"title" validate:"required,min=1,max=255"`
	SortOrder int    `json:"sort_order" validate:"min=0"`
}
