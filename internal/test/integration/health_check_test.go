package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/huypham67/bookmark-service-monolithic/internal/bootstrap"
	healthDTO "github.com/huypham67/bookmark-service-monolithic/internal/dto/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode int
		response   healthDTO.HealthCheckResponse
	}

	testCases := []struct {
		name       string
		appConfig  bootstrap.Config
		setupRedis func(*TestApp)
		expected   expected
	}{
		{
			name: "should return 200 OK with successful health check",
			appConfig: bootstrap.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "instance-1",
			},
			setupRedis: func(app *TestApp) {
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: healthDTO.HealthCheckResponse{
					Message:     "OK",
					ServiceName: "bookmark-service",
					InstanceID:  "instance-1",
				},
			},
		},
		{
			name: "should return 500 when redis connection fails",
			appConfig: bootstrap.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "instance-2",
			},
			setupRedis: func(app *TestApp) {
				app.MockRedis.Close()
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: healthDTO.HealthCheckResponse{
					Message:     "FAILED",
					ServiceName: "bookmark-service",
					InstanceID:  "instance-2",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := setupHealthCheckTestApp(t, tc.appConfig.ServiceName, tc.appConfig.InstanceID)

			tc.setupRedis(app)

			req := httptest.NewRequest(http.MethodGet, "/api/bookmark_service/health-check", nil)
			recorder := httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected.statusCode, recorder.Code)
			assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))

			var actual healthDTO.HealthCheckResponse
			err := json.Unmarshal(recorder.Body.Bytes(), &actual)
			require.NoError(t, err)

			assert.Equal(t, tc.expected.response, actual)
		})
	}
}
