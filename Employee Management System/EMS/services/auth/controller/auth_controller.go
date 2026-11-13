package controller

import (
	"net/http"
	"strings"

	"EMS/services/auth/service"
	apierrors "EMS/shared/errors"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication HTTP requests.
type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new AuthController.
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login handles user login and returns a JWT token.
//
// @Summary      User login
// @Description  Authenticates admin or employee credentials.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body service.LoginRequest true "Login credentials"
// @Success      200 {object} service.LoginResponse
// @Failure      400 {object} apierrors.ErrorResponse
// @Failure      401 {object} apierrors.ErrorResponse
// @Failure      500 {object} apierrors.ErrorResponse
// @Router       /api/auth/login [post]
func (controller *AuthController) Login(c *gin.Context) {
	var request service.LoginRequest

	// Read and validate the login request.
	if err := c.ShouldBindJSON(&request); err != nil {
		apierrors.BadRequest(c, "username and password are required")
		return
	}

	request.Username = strings.TrimSpace(request.Username)

	if request.Username == "" || request.Password == "" {
		apierrors.BadRequest(c, "username and password are required")
		return
	}

	// Authenticate the user and generate a JWT.
	response, err := controller.authService.Login(request)
	if err != nil {
		if strings.Contains(err.Error(), "invalid username or password") {
			apierrors.Unauthorized(c, "invalid username or password")
			return
		}

		apierrors.InternalServerError(c, "login failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"data":    response,
	})
}
