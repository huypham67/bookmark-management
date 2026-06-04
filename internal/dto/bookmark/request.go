package bookmark

type CreateBookmarkRequest struct {
	Description string `json:"description" binding:"required,min=2,max=500"`
	URL         string `json:"url" binding:"required,url"`
}

type ListBookmarksRequest struct {
	Page  int64  `form:"page" binding:"omitempty,min=1"`
	Limit int64  `form:"limit" binding:"omitempty,min=1,max=100"`
	Sort  string `form:"sort" binding:"omitempty,oneof=created_at updated_at code url"`
}
