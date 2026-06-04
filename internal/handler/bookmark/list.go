package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	"github.com/huypham67/bookmark-service/pkg/requestutils"
	"github.com/rs/zerolog/log"
)

// List handles the bookmark list retrieval endpoint with pagination.
//
// @Summary List Bookmarks
// @Description Get paginated list of bookmarks for the authenticated user
// @Tags bookmarks
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number (default: 1)" default(1)
// @Param limit query int false "Items per page, max 100 (default: 10)" default(10)
// @Param sort query string false "Sort field: created_at, updated_at, code, url (default: created_at)" Enums(created_at,updated_at,code,url) default(created_at)
// @Success 200 {object} bookmarkDTO.BookmarkListResponse "List of bookmarks with pagination"
// @Failure 400 {object} gin.H "Invalid query parameters"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/bookmarks [get]
func (h *handler) List(c *gin.Context) {
	userID, err := jwtutils.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	req, err := requestutils.Bind[bookmarkDTO.ListBookmarksRequest](c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters",
		})
		return
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	limit := req.Limit
	if limit == 0 {
		limit = 10
	}

	sort := req.Sort
	if sort == "" {
		sort = "created_at"
	}

	bookmarks, pagination, err := h.service.List(c, userID, page, limit, sort)

	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to list bookmarks")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	bookmarkDataList := make([]bookmarkDTO.BookmarkData, len(bookmarks))
	for i, bm := range bookmarks {
		bookmarkDataList[i] = bookmarkDTO.BookmarkData{
			ID:          bm.ID,
			Code:        bm.Code,
			Description: bm.Description,
			URL:         bm.URL,
			CreatedAt:   bm.CreatedAt,
			UpdatedAt:   bm.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, bookmarkDTO.BookmarkListResponse{
		Data: bookmarkDataList,
		Pagination: bookmarkDTO.Pagination{
			Page:  page,
			Limit: limit,
			Total: pagination.Total,
		},
	})
}
