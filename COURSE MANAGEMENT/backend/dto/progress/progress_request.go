package progress

type UpdateProgressRequest struct {
	WatchedSeconds int  `json:"watched_seconds"`
	Completed      bool `json:"completed"`
}