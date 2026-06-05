package integration

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/huypham67/bookmark-service/pkg/middleware"
	"github.com/stretchr/testify/require"

	"github.com/huypham67/bookmark-service/internal/api"
	authHandler "github.com/huypham67/bookmark-service/internal/handler/auth"
	healthHandler "github.com/huypham67/bookmark-service/internal/handler/health"
	linkHandler "github.com/huypham67/bookmark-service/internal/handler/link"
	profileHandler "github.com/huypham67/bookmark-service/internal/handler/profile"
	linkRepo "github.com/huypham67/bookmark-service/internal/repository/link"
	"github.com/huypham67/bookmark-service/internal/repository/ping"
	userRepo "github.com/huypham67/bookmark-service/internal/repository/user"
	authSvc "github.com/huypham67/bookmark-service/internal/service/auth"
	healthSvc "github.com/huypham67/bookmark-service/internal/service/health"
	linkSvc "github.com/huypham67/bookmark-service/internal/service/link"
	profileSvc "github.com/huypham67/bookmark-service/internal/service/profile"
	"github.com/huypham67/bookmark-service/internal/testutil"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	pkgRedis "github.com/huypham67/bookmark-service/pkg/redis"
	"github.com/huypham67/bookmark-service/pkg/security"
	"github.com/huypham67/bookmark-service/pkg/utils"
)

const (
	testIssuer   = "test-issuer"
	testAudience = "test-audience"
)

// TestApp represents the test application with its dependencies.
type TestApp struct {
	Router    *api.Router
	MockRedis *pkgRedis.MockRedis
}

type AuthenticatedTestApp struct {
	*TestApp
	TokenGenerator jwtutils.TokenGenerator
}

func createTestJWT(t *testing.T) (
	jwtutils.TokenGenerator,
	jwtutils.TokenValidator,
) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(
		rand.Reader,
		2048,
	)
	require.NoError(t, err)

	tokenGenerator, err := jwtutils.NewTokenGenerator(
		privateKey,
		testIssuer,
		testAudience,
		time.Hour,
	)
	require.NoError(t, err)

	tokenValidator, err := jwtutils.NewTokenValidator(
		&privateKey.PublicKey,
		testIssuer,
		testAudience,
	)
	require.NoError(t, err)

	return tokenGenerator, tokenValidator
}

func setupHealthCheckTestApp(t *testing.T, serviceName string, instanceID string) *TestApp {
	t.Helper()

	mockRedis := pkgRedis.NewMockRedis(t)

	pinger := ping.NewRedis(mockRedis.Client)

	healthService := healthSvc.NewService(serviceName, instanceID, pinger)

	healthHandlerInstance := healthHandler.NewHandler(healthService)

	router := api.NewRouter()

	api.RegisterHealthRoutes(
		router.GroupAPI(),
		healthHandlerInstance,
	)

	return &TestApp{
		Router:    router,
		MockRedis: mockRedis,
	}
}

func setupLinkTestApp(t *testing.T) *TestApp {
	t.Helper()

	mockRedis := pkgRedis.NewMockRedis(t)

	linkRepository := linkRepo.NewRepository(mockRedis.Client)

	linkService := linkSvc.NewService(
		linkRepository,
		utils.NewCodeGenerator(),
	)

	linkHandlerInstance := linkHandler.NewHandler(
		linkService,
	)

	router := api.NewRouter()

	api.RegisterLinkRoutes(
		router.GroupV1(),
		linkHandlerInstance,
	)

	return &TestApp{
		Router:    router,
		MockRedis: mockRedis,
	}
}

func createTestTokenGenerator(t *testing.T) jwtutils.TokenGenerator {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	generator, err := jwtutils.NewTokenGenerator(
		privateKey,
		"test-issuer",
		"test-audience",
		time.Hour,
	)
	require.NoError(t, err)

	return generator
}

func setupAuthTestApp(t *testing.T) *TestApp {
	t.Helper()

	mockDB := testutil.NewTestDB(t, &testutil.UserTestDB{})

	userRepository := userRepo.NewRepository(mockDB)

	passwordHasher := security.NewBcryptPasswordHasher()

	tokenGenerator := createTestTokenGenerator(t)

	authService := authSvc.NewService(userRepository, passwordHasher, tokenGenerator)

	authHandlerInstance := authHandler.NewHandler(authService)

	router := api.NewRouter()

	api.RegisterAuthRoutes(
		router.GroupV1(),
		authHandlerInstance,
	)

	return &TestApp{
		Router: router,
	}
}

func setupProfileTestApp(t *testing.T) *AuthenticatedTestApp {
	t.Helper()

	mockDB := testutil.NewTestDB(t, &testutil.UserTestDB{})

	userRepository := userRepo.NewRepository(mockDB)
	profileService := profileSvc.NewService(userRepository)
	profileHandlerInstance := profileHandler.NewHandler(profileService)

	tokenGenerator, tokenValidator := createTestJWT(t)

	router := api.NewRouter()

	jwtMiddleware := middleware.JWTAuth(tokenValidator)

	api.RegisterProfileRoutes(router.GroupV1(), profileHandlerInstance, jwtMiddleware)

	return &AuthenticatedTestApp{
		TestApp: &TestApp{
			Router: router,
		},
		TokenGenerator: tokenGenerator,
	}
}
