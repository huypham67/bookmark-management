package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-service/internal/service/bookmark"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	"github.com/huypham67/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// Delete handles the bookmark deletion endpoint.
//
// @Summary Delete Bookmark
// @Description Delete an existing bookmark
// @Tags bookmarks
// @Produce json
// @Security Bearer
// @Param id path string true "Bookmark ID"
// @Success 200 {object} gin.H "Bookmark deleted successfully"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Bookmark not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/bookmarks/{id} [delete]
func (h *handler) Delete(c *gin.Context) {
	userID, err := jwtutils.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	bookmarkID := c.Param("id")

	if err := h.service.Delete(c, userID, bookmarkID); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("bookmark_id", bookmarkID).
			Msg("failed to delete bookmark")

		switch {
		case errors.Is(err, bookmark.ErrBookmarkNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Bookmark not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.Message("Success"))
}
