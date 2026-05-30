package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetHealthCheck handles the health check endpoint.
//
// @Summary Health Check
// @Description Check application health status and Redis connection
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} github_com_huypham67_bookmark_service_internal_dto_response.HealthCheckResponse
// @Failure 500 {object} github_com_huypham67_bookmark_service_internal_dto_response.HealthCheckResponse
// @Router /health-check [get]
func (h *handler) GetHealthCheck(c *gin.Context) {
	res := h.service.GetStatus(c)
	if res.Message == "FAILED" {
		log.Error().
			Str("message", res.Message).
			Msg("500 - health check failed")

		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, res)
}
