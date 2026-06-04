package bookmark

import "time"

type BookmarkData struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type BookmarkResponse struct {
	Data    *BookmarkData `json:"data"`
	Message string        `json:"message"`
}
