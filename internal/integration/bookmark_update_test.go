package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBookmarkEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		bookmarkID  string
		requestBody string
		setupAuth   func(t *testing.T, app *AuthenticatedTestApp, req *http.Request)
		expected    expected
	}{
		{
			name:       "should return 200 when bookmark is updated successfully",
			bookmarkID: "bookmark-1-1",
			requestBody: `{
				"description": "Updated Description",
				"url": "https://updated-example.com"
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
				statusCode:   http.StatusOK,
				bodyContains: "Success",
			},
		},
		{
			name:        "should return 400 when request body is invalid JSON",
			bookmarkID:  "bookmark-1-1",
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
				bodyContains: "Invalid request",
			},
		},
		{
			name:       "should return 401 when authorization header is missing",
			bookmarkID: "bookmark-1-1",
			requestBody: `{
				"description": "Updated Description",
				"url": "https://updated-example.com"
			}`,
			setupAuth: func(t *testing.T, app *AuthenticatedTestApp, req *http.Request) {
				// No auth header set
			},
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "missing authorization header",
			},
		},
		{
			name:       "should return 404 when bookmark does not exist",
			bookmarkID: "nonexistent-bookmark",
			requestBody: `{
				"description": "Updated Description",
				"url": "https://updated-example.com"
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
				statusCode:   http.StatusNotFound,
				bodyContains: "Bookmark not found",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := setupBookmarkTestApp(t)

			httpRequest := httptest.NewRequest(
				http.MethodPut,
				"/api/bookmark_service/v1/bookmarks/"+tc.bookmarkID,
				bytes.NewBufferString(tc.requestBody),
			)

			httpRequest.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, app, httpRequest)

			httpRecorder := httptest.NewRecorder()

			app.Router.ServeHTTP(httpRecorder, httpRequest)

			assert.Equal(t, tc.expected.statusCode, httpRecorder.Code)
			assert.Contains(t, httpRecorder.Body.String(), tc.expected.bodyContains)
		})
	}
}
