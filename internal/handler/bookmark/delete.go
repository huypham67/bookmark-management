package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/huypham67/bookmark-service/internal/service/bookmark"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	"github.com/huypham67/bookmark-service/pkg/requestutils"
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
// @Failure 400 {object} gin.H "Invalid bookmark ID"
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

	req, err := requestutils.Bind[bookmarkDTO.DeleteBookmarkRequest](c)

	if err != nil {
		log.Warn().
			Err(err).
			Msg("invalid delete bookmark request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := h.service.Delete(c, userID, req.ID); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("bookmark_id", req.ID).
			Msg("failed to delete bookmark")

		switch {
		case errors.Is(err, bookmark.ErrBookmarkNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Bookmark not found",
			})
		case errors.Is(err, bookmark.ErrBadRequest):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
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
