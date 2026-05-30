package link

import (
	"net/http"

	"github.com/gin-gonic/gin"
	linkDTO "github.com/huypham67/bookmark-service/internal/dto/link"
	"github.com/rs/zerolog/log"
)

// ShortenURL handles the URL shortening endpoint.
//
// @Summary Shorten URL
// @Description Create a shortened URL code and save it to Redis
// @Tags links
// @Accept json
// @Produce json
// @Param request body request.ShortenURLRequest true "URL to shorten"
// @Success 200 {object} github_com_huypham67_bookmark_service_internal_dto_response.ShortenURLResponse "Shorten URL generated successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/links/shorten [post]
func (h *handler) ShortenURL(c *gin.Context) {
	var req linkDTO.ShortenURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	code, err := h.service.ShortenURL(c, req)

	if err != nil {
		log.Error().
			Err(err).
			Str("url", req.Url).
			Int64("exp", req.Exp).
			Msg("500 - failed to shorten URL")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})

		return
	}

	c.JSON(http.StatusOK, linkDTO.ShortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully",
	})
}
