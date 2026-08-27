package lesson

type UpdateLessonRequest struct {
	Title           string `json:"title" validate:"required,min=1,max=255"`
	VideoURL        string `json:"video_url" validate:"required,url"`
	DurationSeconds int    `json:"duration_seconds" validate:"min=0"`
	IsPreview       bool   `json:"is_preview"`
	SortOrder       int    `json:"sort_order" validate:"min=0"`
}
