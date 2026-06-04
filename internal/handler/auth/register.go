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

// Register handles the user registration endpoint.
//
// @Summary Register User
// @Description Register a new user with email, username, and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body authDTO.RegisterUserRequest true "User registration data"
// @Success 201 {object} authDTO.RegisterUserResponse "User registered successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 409 {object} gin.H "User already exists"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/users/register [post]
func (h *handler) Register(c *gin.Context) {
	req, err := requestutils.Bind[authDTO.RegisterUserRequest](c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	user, err := h.service.RegisterUser(c, *req)

	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Str("username", req.Username).
			Msg("failed to register user")

		switch {
		case errors.Is(err, auth.ErrEmailAlreadyRegistered), errors.Is(err, auth.ErrUsernameAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, authDTO.RegisterUserResponse{
		Data: authDTO.UserData{
			ID:          user.ID,
			DisplayName: user.DisplayName,
			Username:    user.Username,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
		},
		Message: "Register an user successfully!",
	})
}
