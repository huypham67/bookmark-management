package link

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	code := c.Param("code")

	url, err := h.service.GetOriginalURL(c, code)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short link not found",
		})
		return
	}

	c.Redirect(http.StatusFound, url)
}
