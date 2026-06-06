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

// Update handles the bookmark update endpoint.
//
// @Summary Update Bookmark
// @Description Update an existing bookmark (description and/or URL)
// @Tags bookmarks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Bookmark ID"
// @Param request body bookmarkDTO.UpdateBookmarkRequest true "Updated bookmark data"
// @Success 200 {object} gin.H "Bookmark updated successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Bookmark not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/bookmarks/{id} [put]
func (h *handler) Update(c *gin.Context) {
	userID, err := jwtutils.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	bookmarkID := c.Param("id")

	req, err := requestutils.Bind[bookmarkDTO.UpdateBookmarkRequest](c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.Update(c, userID, bookmarkID, *req); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("bookmark_id", bookmarkID).
			Msg("failed to update bookmark")

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
