package health

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-service/internal/service/health"
)

type Handler interface {
	GetHealthCheck(c *gin.Context)
}

type handler struct {
	service health.Service
}

func NewHandler(service health.Service) Handler {
	return &handler{
		service: service,
	}
}
