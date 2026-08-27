package filematerial

type FileMaterialResponse struct {
	ID       string `json:"id"`
	LessonID string `json:"lesson_id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}