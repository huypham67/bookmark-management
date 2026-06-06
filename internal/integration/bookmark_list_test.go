package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListBookmarksEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode    int
		bodyContains  string
		expectedCount int
		hasMoreData   bool
	}

	testCases := []struct {
		name        string
		queryParams string
		setupAuth   func(t *testing.T, app *AuthenticatedTestApp, req *http.Request)
		expected    expected
	}{
		{
			name:        "should return 200 with bookmarks list",
			queryParams: "?page=1&limit=10",
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
				statusCode:    http.StatusOK,
				bodyContains:  "Bookmarks retrieved successfully!",
				expectedCount: 8,
				hasMoreData:   false,
			},
		},
		{
			name:        "should return 401 when authorization header is missing",
			queryParams: "?page=1&limit=10",
			setupAuth: func(t *testing.T, app *AuthenticatedTestApp, req *http.Request) {
				// No auth header set
			},
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "missing authorization header",
			},
		},
		{
			name:        "should return 200 with empty bookmarks list when user has no bookmarks",
			queryParams: "?page=10&limit=10",
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
				statusCode:    http.StatusOK,
				bodyContains:  "Bookmarks retrieved successfully!",
				expectedCount: 0,
				hasMoreData:   false,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := setupBookmarkTestApp(t)

			httpRequest := httptest.NewRequest(
				http.MethodGet,
				"/api/bookmark_service/v1/bookmarks"+tc.queryParams,
				nil,
			)

			tc.setupAuth(t, app, httpRequest)

			httpRecorder := httptest.NewRecorder()

			app.Router.ServeHTTP(httpRecorder, httpRequest)

			assert.Equal(t, tc.expected.statusCode, httpRecorder.Code)
			assert.Contains(t, httpRecorder.Body.String(), tc.expected.bodyContains)

			if tc.expected.statusCode == http.StatusOK {
				var resp bookmarkDTO.BookmarkListResponse
				err := json.Unmarshal(httpRecorder.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, tc.expected.expectedCount, len(resp.Data))
				assert.NotEmpty(t, resp.Pagination)
				assert.Equal(t, int64(10), resp.Pagination.Limit)
			}
		})
	}
}
