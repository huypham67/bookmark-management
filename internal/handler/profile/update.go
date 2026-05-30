package profile

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	profileDTO "github.com/huypham67/bookmark-service/internal/dto/profile"
	"github.com/huypham67/bookmark-service/internal/service/profile"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	"github.com/huypham67/bookmark-service/pkg/requestutils"
	"github.com/rs/zerolog/log"
)

// UpdateUserInfo handles the user info update endpoint.
//
// @Summary Update User Info
// @Description Update authenticated user's display name and email
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body github_com_huypham67_bookmark_service_internal_dto_request.UpdateUserRequest true "User update data"
// @Success 200 {object} github_com_huypham67_bookmark_service_internal_dto_response.UpdateUserResponse "User updated successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 409 {object} gin.H "Email already exists"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/self/info [put]
func (h *handler) UpdateUserInfo(c *gin.Context) {
	userID, err := jwtutils.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	req, err := requestutils.Bind[profileDTO.UpdateUserRequest](c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdateUserInfo(c, userID, *req); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("email", req.Email).
			Msg("failed to update user info")

		// Check if email already exists
		if errors.Is(err, profile.ErrEmailAlreadyRegistered) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})

		return
	}

	c.JSON(http.StatusOK, profileDTO.UpdateUserResponse{
		Message: "Edit current user successfully!",
	})
}
