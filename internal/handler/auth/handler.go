package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-service/internal/service/auth"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type handler struct {
	service auth.Service
}

func NewHandler(service auth.Service) Handler {
	return &handler{
		service: service,
	}
}
