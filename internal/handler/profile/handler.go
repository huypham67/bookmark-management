package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-service/internal/service/profile"
)

type Handler interface {
	GetUserInfo(c *gin.Context)
	UpdateUserInfo(c *gin.Context)
}

type handler struct {
	service profile.Service
}

func NewHandler(service profile.Service) Handler {
	return &handler{
		service: service,
	}
}
