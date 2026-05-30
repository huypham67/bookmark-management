package link

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-service/internal/service/link"
)

type Handler interface {
	ShortenURL(c *gin.Context)
	RedirectToURL(c *gin.Context)
}

type handler struct {
	service link.Service
}

func NewHandler(service link.Service) Handler {
	return &handler{
		service: service,
	}
}
