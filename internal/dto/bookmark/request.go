package bookmark

type CreateBookmarkRequest struct {
	Description string `json:"description" binding:"required,min=2,max=500"`
	URL         string `json:"url" binding:"required,url"`
}
