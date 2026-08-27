package course

type CreateCourseRequest struct {
	CategoryID      string  `json:"category_id" validate:"required,uuid"`
	Title           string  `json:"title" validate:"required,max=255"`
	Description     string  `json:"description"`
	ThumbnailURL    *string `json:"thumbnail_url"`
	PreviewVideoURL *string `json:"preview_video_url"`
	Price           float64 `json:"price" validate:"gte=0"`
}
