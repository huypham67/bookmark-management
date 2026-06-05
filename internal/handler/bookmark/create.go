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

// Create handles the bookmark creation endpoint.
//
// @Summary Create Bookmark
// @Description Create a new bookmark for the authenticated user
// @Tags bookmarks
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body bookmarkDTO.CreateBookmarkRequest true "Bookmark creation data"
// @Success 201 {object} bookmarkDTO.BookmarkResponse "Bookmark created successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/bookmarks [post]
func (h *handler) Create(c *gin.Context) {
	userID, err := jwtutils.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	req, err := requestutils.Bind[bookmarkDTO.CreateBookmarkRequest](c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	bm, err := h.service.Create(c, userID, *req)

	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("url", req.URL).
			Msg("failed to create bookmark")

		switch {
		case errors.Is(err, bookmark.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, response.Success(
		&bookmarkDTO.BookmarkData{
			ID:          bm.ID,
			Code:        bm.Code,
			Description: bm.Description,
			URL:         bm.URL,
			CreatedAt:   bm.CreatedAt,
			UpdatedAt:   bm.UpdatedAt,
		},
		"Bookmark created successfully!",
	))
}
