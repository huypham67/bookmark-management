package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-service-monolithic/middleware"
	pkgRedis "github.com/huypham67/bookmark-service-monolithic/pkg/redis"
	"github.com/huypham67/bookmark-service-monolithic/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	authHandler "github.com/huypham67/bookmark-service-monolithic/internal/handler/auth"
	bookmarkHandler "github.com/huypham67/bookmark-service-monolithic/internal/handler/bookmark"
	healthHandler "github.com/huypham67/bookmark-service-monolithic/internal/handler/health"
	linkHandler "github.com/huypham67/bookmark-service-monolithic/internal/handler/link"
	profileHandler "github.com/huypham67/bookmark-service-monolithic/internal/handler/profile"
	bookmarkRepo "github.com/huypham67/bookmark-service-monolithic/internal/repository/bookmark"
	cacheRepo "github.com/huypham67/bookmark-service-monolithic/internal/repository/cache"
	linkRepo "github.com/huypham67/bookmark-service-monolithic/internal/repository/link"
	"github.com/huypham67/bookmark-service-monolithic/internal/repository/ping"
	"github.com/huypham67/bookmark-service-monolithic/internal/repository/user"
	authSvc "github.com/huypham67/bookmark-service-monolithic/internal/service/auth"
	bookmarkSvc "github.com/huypham67/bookmark-service-monolithic/internal/service/bookmark"
	bookmarkCacheSvc "github.com/huypham67/bookmark-service-monolithic/internal/service/bookmark/cache"
	healthSvc "github.com/huypham67/bookmark-service-monolithic/internal/service/health"
	linkSvc "github.com/huypham67/bookmark-service-monolithic/internal/service/link"
	profileSvc "github.com/huypham67/bookmark-service-monolithic/internal/service/profile"
	"github.com/huypham67/bookmark-service-monolithic/pkg/jwt"
	"github.com/huypham67/bookmark-service-monolithic/pkg/jwtprovider"
	"github.com/huypham67/bookmark-service-monolithic/pkg/security"
	"github.com/huypham67/bookmark-service-monolithic/pkg/utils"
)

// Container holds the application's dependencies and initialized services.
// It serves as the single source of truth for all infrastructure and business logic components.
type Container struct {
	// Infrastructure
	Config    *Config
	DB        *gorm.DB
	Redis     *redis.Client
	CacheRepo cacheRepo.Repository

	// Handlers
	HealthHandler   healthHandler.Handler
	AuthHandler     authHandler.Handler
	ProfileHandler  profileHandler.Handler
	LinkHandler     linkHandler.Handler
	BookmarkHandler bookmarkHandler.Handler

	// Middleware
	JWTMiddleware gin.HandlerFunc
}

// NewContainer initializes the application container by loading configuration,
// setting up infrastructure clients, and initializing all handlers with their dependencies.
func NewContainer() (*Container, error) {
	cfg, err := NewConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		return nil, err
	}

	db, err := sqldb.NewDBClient("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize postgres client")
		return nil, err
	}

	db, err = sqldb.RunMigration(db, "migrations")
	if err != nil {
		log.Error().Err(err).Msg("failed to run database migrations")
		return nil, err
	}

	rdb, err := pkgRedis.NewRedisClient("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize redis client")
		return nil, err
	}

	jwtProvider, err := jwtprovider.NewProvider("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize jwt provider")
		return nil, err
	}

	jwtMiddleware := middleware.JWTAuth(jwtProvider.Validator())

	// Initialize shared infrastructure
	cacheRepository := cacheRepo.NewRedis(rdb)

	healthHandlerInstance := initHealthHandler(cfg, rdb)
	authHandlerInstance, err := initAuthHandler(db, jwtProvider.Generator())
	if err != nil {
		return nil, err
	}
	profileHandlerInstance := initProfileHandler(db)
	linkHandlerInstance := initLinkHandler(rdb, db)
	bookmarkHandlerInstance := initBookmarkHandler(db, cacheRepository)

	return &Container{
		Config:          cfg,
		DB:              db,
		Redis:           rdb,
		CacheRepo:       cacheRepository,
		HealthHandler:   healthHandlerInstance,
		AuthHandler:     authHandlerInstance,
		ProfileHandler:  profileHandlerInstance,
		LinkHandler:     linkHandlerInstance,
		BookmarkHandler: bookmarkHandlerInstance,
		JWTMiddleware:   jwtMiddleware,
	}, nil
}

func initAuthHandler(db *gorm.DB, tokenGenerator jwt.TokenGenerator) (authHandler.Handler, error) {
	userRepository := user.NewRepository(db)
	passwordHasher := security.NewBcryptPasswordHasher()

	authService := authSvc.NewService(userRepository, passwordHasher, tokenGenerator)
	return authHandler.NewHandler(authService), nil
}

func initProfileHandler(db *gorm.DB) profileHandler.Handler {
	userRepository := user.NewRepository(db)
	service := profileSvc.NewService(userRepository)
	return profileHandler.NewHandler(service)
}

func initHealthHandler(cfg *Config, redisClient *redis.Client) healthHandler.Handler {
	pinger := ping.NewRedis(redisClient)
	healthService := healthSvc.NewService(cfg.ServiceName, cfg.InstanceID, pinger)
	return healthHandler.NewHandler(healthService)
}

func initLinkHandler(redisClient *redis.Client, db *gorm.DB) linkHandler.Handler {
	linkRepository := linkRepo.NewRepository(redisClient)
	codeGenerator := utils.NewCodeGenerator()
	bookmarkResolver := bookmarkRepo.NewRepository(db)
	service := linkSvc.NewService(linkRepository, codeGenerator, bookmarkResolver)
	return linkHandler.NewHandler(service)
}

func initBookmarkHandler(db *gorm.DB, cacheRepository cacheRepo.Repository) bookmarkHandler.Handler {
	bookmarkRepository := bookmarkRepo.NewRepository(db)
	bookmarkService := bookmarkSvc.NewService(bookmarkRepository)

	cacheService := bookmarkCacheSvc.NewBookmarkService(bookmarkService, cacheRepository)

	return bookmarkHandler.NewHandler(cacheService)
}

// Close gracefully shuts down the database and Redis clients, ensuring that all resources are properly released.
func (c *Container) Close() error {
	if c.Redis != nil {
		_ = c.Redis.Close()
	}

	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	return nil
}
