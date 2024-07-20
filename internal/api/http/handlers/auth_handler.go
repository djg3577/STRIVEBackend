package handlers

import (
	"STRIVEBackend/internal/service"
	"STRIVEBackend/internal/util"
	"STRIVEBackend/pkg/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

func (h *AuthHandler) DecodeJWT(c *gin.Context) {
	user, err := h.Service.AuthenticateUser(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR DECODING JWT" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := h.Service.AuthenticateUser(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userID", user.ID)

		c.Next()
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := util.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields of: " + err.Error()})
		return
	}

	if h.Service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service is nil"})
		return
	}

	userID, err := h.Service.SignUp(&user)
	if err != nil {
		handleSignUpError(c, err)
		return
	}

	token, err := util.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "token": token})
}

func handleSignUpError(c *gin.Context, err error) {
	switch {
	case strings.Contains(err.Error(), "users_username_key"):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists. Please pick a new one"})
	case strings.Contains(err.Error(), "users_email_key"):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already associated with an account. Please login instead"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
	}
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Code  int    `json:"code"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.VerifyEmail(req.Email, req.Code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}


func (h *AuthHandler) GitHubAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
					c.Abort()
					return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
					c.Abort()
					return
			}

			token := bearerToken[1]
			githubUser, err := h.Service.GetGitHubUser(token)
			if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid GitHub token"})
					c.Abort()
					return
			}

			internalUserId, err := h.Service.GetOrCreateUserIdFromGithub(githubUser.ID)
			if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create user from GitHub"})
					c.Abort()
					return
			}

			// Store the GitHub user ID in the context
			c.Set("githubUserId", githubUser.ID)
			c.Set("userID", internalUserId)
			c.Next()
	}
}

func (h *AuthHandler) GitHubLogin(c *gin.Context){
	var request struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.Service.ExchangeGitHubCode(request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error github auth failed:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}