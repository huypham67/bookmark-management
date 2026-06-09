package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteBookmarkEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		bookmarkID  string
		setupAuth   func(t *testing.T, app *AuthenticatedTestApp, req *http.Request)
		setupCache  func(t *testing.T, app *AuthenticatedTestApp)
		verifyCache func(t *testing.T, app *AuthenticatedTestApp)
		expected    expected
	}{
		{
			name:       "should return 401 when authorization header is missing",
			bookmarkID: "bookmark-1-1",
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
		{
			name:       "should return 200 and invalidate cache when bookmark is deleted successfully",
			bookmarkID: "bookmark-1-1",
			setupCache: func(t *testing.T, app *AuthenticatedTestApp) {
				seedBookmarkListCache(t, app, 1, 10, "created_at")
				require.True(t, app.MockRedis.Server.Exists(bookmarkCacheHashKey(cacheSeedUserID)))
			},
			setupAuth: func(t *testing.T, app *AuthenticatedTestApp, req *http.Request) {
				token, err := app.TokenGenerator.GenerateToken(
					"user-uuid-1",
					"testuser1",
					"testuser1@gmail.com",
				)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			verifyCache: func(t *testing.T, app *AuthenticatedTestApp) {
				assert.False(t, app.MockRedis.Server.Exists(bookmarkCacheHashKey(cacheSeedUserID)))
			},
			expected: expected{
				statusCode:   http.StatusOK,
				bodyContains: "Success",
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := setupBookmarkTestApp(t)

			if tc.setupCache != nil {
				tc.setupCache(t, app)
			}

			httpRequest := httptest.NewRequest(
				http.MethodDelete,
				"/api/bookmark_service/v1/bookmarks/"+tc.bookmarkID,
				nil,
			)

			tc.setupAuth(t, app, httpRequest)

			httpRecorder := httptest.NewRecorder()

			app.Router.ServeHTTP(httpRecorder, httpRequest)

			assert.Equal(t, tc.expected.statusCode, httpRecorder.Code)
			assert.Contains(t, httpRecorder.Body.String(), tc.expected.bodyContains)

			if tc.verifyCache != nil {
				tc.verifyCache(t, app)
			}
		})
	}
}
