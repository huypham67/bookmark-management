package link

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	linkDTO "github.com/huypham67/bookmark-service/internal/dto/link"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
	"github.com/huypham67/bookmark-service/pkg/requestutils"
	"github.com/huypham67/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// RedirectToURL handles the redirect endpoint.
//
// @Summary Redirect to Original URL
// @Description Redirect user to the original URL based on the shortened code
// @Tags links
// @Accept json
// @Produce json
// @Param code path string true "Shortened code"
// @Success 301 "Redirect successful"
// @Failure 404 {object} gin.H "Short link not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/links/redirect/{code} [get]
func (h *handler) RedirectToURL(c *gin.Context) {
	req, err := requestutils.Bind[linkDTO.RedirectRequest](c)

	if err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	code := req.Code

	url, err := h.service.GetOriginalURL(c, code)

	if err != nil {
		if errors.Is(err, dbutils.ErrRecordNotFoundType) {
			response.NotFound(c, "Short link not found")
			return
		}

		log.Error().
			Err(err).
			Str("code", code).
			Msg("failed to retrieve original URL")

		response.InternalServerError(c)
		return
	}

	c.Redirect(http.StatusFound, url)
}
