package bootstrap

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-service/pkg/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/huypham67/bookmark-service/docs"
	"github.com/huypham67/bookmark-service/internal/api"
	authHandler "github.com/huypham67/bookmark-service/internal/handler/auth"
	bookmarkHandler "github.com/huypham67/bookmark-service/internal/handler/bookmark"
	healthHandler "github.com/huypham67/bookmark-service/internal/handler/health"
	linkHandler "github.com/huypham67/bookmark-service/internal/handler/link"
	profileHandler "github.com/huypham67/bookmark-service/internal/handler/profile"
	bookmarkRepo "github.com/huypham67/bookmark-service/internal/repository/bookmark"
	linkRepo "github.com/huypham67/bookmark-service/internal/repository/link"
	"github.com/huypham67/bookmark-service/internal/repository/ping"
	"github.com/huypham67/bookmark-service/internal/repository/user"
	authSvc "github.com/huypham67/bookmark-service/internal/service/auth"
	bookmarkSvc "github.com/huypham67/bookmark-service/internal/service/bookmark"
	healthSvc "github.com/huypham67/bookmark-service/internal/service/health"
	linkSvc "github.com/huypham67/bookmark-service/internal/service/link"
	profileSvc "github.com/huypham67/bookmark-service/internal/service/profile"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	"github.com/huypham67/bookmark-service/pkg/logger"
	pkgRedis "github.com/huypham67/bookmark-service/pkg/redis"
	"github.com/huypham67/bookmark-service/pkg/security"
	"github.com/huypham67/bookmark-service/pkg/sqldb"
	"github.com/huypham67/bookmark-service/pkg/utils"
)

// App represents application runtime container.
type App struct {
	config      *Config
	router      *api.Router
	redisClient *redis.Client
	dbClient    *gorm.DB
}

// NewApp initializes application dependencies.
func NewApp() (*App, error) {
	if err := logger.NewLoggerClient(""); err != nil {
		return nil, err
	}

	cfg, err := NewConfig()

	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to load config")

		return nil, err
	}

	setupSwaggerConfig(cfg)

	redisClient, err := initRedisClient()

	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to initialize redis client")

		return nil, err
	}

	dbClient, err := initPostgresClient()

	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to initialize postgres client")

		return nil, err
	}

	router := api.NewRouter()

	if err := registerRoutes(router, cfg, redisClient, dbClient); err != nil {
		return nil, err
	}

	log.Info().
		Msg("application initialized successfully")

	return &App{
		config:      cfg,
		router:      router,
		redisClient: redisClient,
	}, nil
}

func registerRoutes(router *api.Router, cfg *Config, redisClient *redis.Client, dbClient *gorm.DB) error {
	apiGroup := router.GroupAPI()
	apiV1Group := router.GroupV1()

	healthHandlerInstance := initHealthHandler(cfg, redisClient)
	linkHandlerInstance := initLinkHandler(redisClient)
	authHandlerInstance, err := initAuthHandler(dbClient)
	if err != nil {
		return err
	}
	profileHandlerInstance := initProfileHandler(dbClient)
	bookmarkHandlerInstance := initBookmarkHandler(dbClient)

	jwtMiddleware := initJWTMiddleware()

	api.RegisterHealthRoutes(apiGroup, healthHandlerInstance)
	api.RegisterLinkRoutes(apiV1Group, linkHandlerInstance)
	api.RegisterAuthRoutes(apiV1Group, authHandlerInstance)
	api.RegisterProfileRoutes(apiV1Group, profileHandlerInstance, jwtMiddleware)
	api.RegisterBookmarkRoutes(apiV1Group, bookmarkHandlerInstance, jwtMiddleware)

	return nil
}

func initRedisClient() (*redis.Client, error) {
	return pkgRedis.NewRedisClient("")
}

func initPostgresClient() (*gorm.DB, error) {
	return sqldb.NewDBClient("")
}

const tokenIssuer = "bookmark-service"
const tokenAudience = "bookmark-service"
const tokenExpiration = time.Hour * 1
const privateKeyPath = "keys/private.pem"
const publicKeyPath = "keys/public.pem"

func initJWTMiddleware() gin.HandlerFunc {
	// Load public key for JWT token validation
	publicKey, err := jwtutils.LoadRSAPublicKeyFromFile(publicKeyPath)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to load public key for JWT")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}

	tokenValidator, err := jwtutils.NewTokenValidator(publicKey, tokenIssuer, tokenAudience)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to create token validator")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}

	return middleware.JWTAuth(tokenValidator)
}

func initAuthHandler(db *gorm.DB) (authHandler.Handler, error) {
	userRepository := user.NewRepository(db)
	passwordHasher := security.NewBcryptPasswordHasher()

	// Load private key for JWT token generation
	privateKey, err := jwtutils.LoadRSAPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to load private key for JWT")
		return nil, err
	}

	tokenGenerator, err := jwtutils.NewTokenGenerator(privateKey, tokenIssuer, tokenAudience, tokenExpiration)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to create token generator")
		return nil, err
	}

	// Create auth service (for registration and login)
	authService := authSvc.NewService(userRepository, passwordHasher, tokenGenerator)
	return authHandler.NewHandler(authService), nil
}

func initProfileHandler(db *gorm.DB) profileHandler.Handler {
	userRepository := user.NewRepository(db)

	// Create profile service (for getting and updating user info)
	service := profileSvc.NewService(userRepository)
	return profileHandler.NewHandler(service)
}

func initHealthHandler(cfg *Config, redisClient *redis.Client) healthHandler.Handler {
	pinger := ping.NewRedis(redisClient)

	healthService := healthSvc.NewService(cfg.ServiceName, cfg.InstanceID, pinger)

	return healthHandler.NewHandler(healthService)
}

func initLinkHandler(redisClient *redis.Client) linkHandler.Handler {
	linkRepository := linkRepo.NewRepository(redisClient)

	codeGenerator := utils.NewCodeGenerator()

	service := linkSvc.NewService(linkRepository, codeGenerator)

	return linkHandler.NewHandler(service)
}

func initBookmarkHandler(db *gorm.DB) bookmarkHandler.Handler {
	bookmarkRepository := bookmarkRepo.NewRepository(db)
	codeGenerator := utils.NewCodeGenerator()

	service := bookmarkSvc.NewService(bookmarkRepository, codeGenerator)
	return bookmarkHandler.NewHandler(service)
}

func setupSwaggerConfig(cfg *Config) {
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = getSwaggerSchemes(cfg)
	docs.SwaggerInfo.BasePath = cfg.HostName
}

func getSwaggerSchemes(cfg *Config) []string {
	if cfg.SwaggerSchemes != "" {
		return parseSchemes(cfg.SwaggerSchemes)
	}

	if cfg.Environment == "production" {
		return []string{"https"}
	}

	return []string{"http"}
}

func parseSchemes(schemesStr string) []string {
	schemes := make([]string, 0)
	for _, scheme := range strings.Split(schemesStr, ",") {
		if trimmed := strings.TrimSpace(scheme); trimmed != "" {
			schemes = append(schemes, trimmed)
		}
	}
	return schemes
}

// Run starts HTTP server.
func (a *App) Run() error {
	return http.ListenAndServe(
		":"+a.config.AppPort,
		a.router,
	)
}

// Close closes all application resources.
func (a *App) Close() error {
	if a.redisClient != nil {
		_ = a.redisClient.Close()
	}

	if a.dbClient != nil {
		sqlDB, err := a.dbClient.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	return nil
}
