package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	authDTO "github.com/huypham67/bookmark-service/internal/dto/auth"
	"github.com/huypham67/bookmark-service/internal/service/auth"
	"github.com/huypham67/bookmark-service/pkg/requestutils"
	"github.com/rs/zerolog/log"
)

// Login handles the user login endpoint.
//
// @Summary Login User
// @Description Login a user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body github_com_huypham67_bookmark_service_internal_dto_request.LoginRequest true "User login data"
// @Success 200 {object} github_com_huypham67_bookmark_service_internal_dto_response.LoginResponse "User logged in successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Invalid credentials"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/users/login [post]
func (h *handler) Login(c *gin.Context) {
	req, err := requestutils.Bind[authDTO.LoginRequest](c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	token, err := h.service.LoginUser(c, *req)

	if err != nil {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Msg("failed to login user")

		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})

		return
	}

	c.JSON(http.StatusOK, authDTO.LoginResponse{
		Data:    token,
		Message: "Logged in successfully!",
	})
}
