package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBookmarkEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		requestBody string
		setupAuth   func(t *testing.T, app *AuthenticatedTestApp, req *http.Request)
		expected    expected
	}{
		{
			name: "should return 201 when bookmark is created successfully",
			requestBody: `{
				"description": "New Bookmark",
				"url": "https://new-example.com"
			}`,
			setupAuth: func(t *testing.T, app *AuthenticatedTestApp, req *http.Request) {
				token, err := app.TokenGenerator.GenerateToken(
					"user-uuid-1",
					"testuser1",
					"testuser1@gmail.com",
				)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expected: expected{
				statusCode:   http.StatusCreated,
				bodyContains: "Bookmark created successfully!",
			},
		},
		{
			name: "should return 400 when user does not exist",
			requestBody: `{
				"description": "Test Bookmark",
				"url": "https://test-example.com"
			}`,
			setupAuth: func(t *testing.T, app *AuthenticatedTestApp, req *http.Request) {
				token, err := app.TokenGenerator.GenerateToken(
					"nonexistent-user-id",
					"nonexistent-user",
					"nonexistent@example.com",
				)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expected: expected{
				statusCode:   http.StatusBadRequest,
				bodyContains: "Invalid bookmark request",
			},
		},
		{
			name:        "should return 400 when request body is invalid JSON",
			requestBody: `{invalid json}`,
			setupAuth: func(t *testing.T, app *AuthenticatedTestApp, req *http.Request) {
				token, err := app.TokenGenerator.GenerateToken(
					"user-uuid-1",
					"testuser1",
					"testuser1@gmail.com",
				)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expected: expected{
				statusCode:   http.StatusBadRequest,
				bodyContains: "Invalid request body",
			},
		},
		{
			name: "should return 401 when authorization header is missing",
			requestBody: `{
				"description": "Test Bookmark",
				"url": "https://test-example.com"
			}`,
			setupAuth: func(t *testing.T, app *AuthenticatedTestApp, req *http.Request) {
				// No auth header set
			},
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "missing authorization header",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := setupBookmarkTestApp(t)

			httpRequest := httptest.NewRequest(
				http.MethodPost,
				"/api/bookmark_service/v1/bookmarks",
				bytes.NewBufferString(tc.requestBody),
			)

			httpRequest.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, app, httpRequest)

			httpRecorder := httptest.NewRecorder()

			app.Router.ServeHTTP(httpRecorder, httpRequest)

			assert.Equal(t, tc.expected.statusCode, httpRecorder.Code)
			assert.Contains(t, httpRecorder.Body.String(), tc.expected.bodyContains)

			if tc.expected.statusCode == http.StatusCreated {
				var resp bookmarkDTO.BookmarkResponse
				err := json.Unmarshal(httpRecorder.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Data.ID)
				assert.NotEmpty(t, resp.Data.Code)
				assert.NotEmpty(t, resp.Data.Description)
				assert.NotEmpty(t, resp.Data.URL)
			}
		})
	}
}
