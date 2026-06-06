package bootstrap

import (
	"github.com/huypham67/bookmark-service/internal/api"
)

// SetupRoutes registers all API routes with the provided router and dependency container.
func SetupRoutes(router *api.Router, container *Container) {
	apiGroup := router.GroupAPI()
	apiV1Group := router.GroupV1()

	// Register all feature routes in order
	api.RegisterHealthRoutes(apiGroup, container.HealthHandler)
	api.RegisterAuthRoutes(apiV1Group, container.AuthHandler)
	api.RegisterProfileRoutes(apiV1Group, container.ProfileHandler, container.JWTMiddleware)
	api.RegisterLinkRoutes(apiV1Group, container.LinkHandler)
	api.RegisterBookmarkRoutes(apiV1Group, container.BookmarkHandler, container.JWTMiddleware)
}
