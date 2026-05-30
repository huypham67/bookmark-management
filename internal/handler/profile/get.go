package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	profileDTO "github.com/huypham67/bookmark-service/internal/dto/profile"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	"github.com/rs/zerolog/log"
)

// GetUserInfo handles the user info endpoint.
//
// @Summary Get User Info
// @Description Get authenticated user information from JWT token
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} github_com_huypham67_bookmark_service_internal_dto_response.UserResponse "User information"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/self/info [get]
func (h *handler) GetUserInfo(c *gin.Context) {
	userID, err := jwtutils.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, err := h.service.GetUserInfo(c, userID)

	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to get user info")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})

		return
	}

	c.JSON(http.StatusOK, profileDTO.UserResponse{
		Data: &profileDTO.UserData{
			ID:          user.ID,
			DisplayName: user.DisplayName,
			Username:    user.Username,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
		},
		Message: "User information retrieved successfully!",
	})
}
