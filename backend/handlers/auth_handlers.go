package handlers

import (
	"mis-system/database"
	"mis-system/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	jwtKey             = []byte("your_secret_key")                           // In production, this should be an environment variable
	accessTokenExp     = 15 * time.Minute                                    // 15 minutes
	refreshTokenExp    = 30 * 24 * time.Hour                                 // 30 days
	googleClientID     = "your-google-client-id"                             // In production, this should be an environment variable
	googleClientSecret = "your-google-client-secret"                         // In production, this should be an environment variable
	googleRedirectURL  = "http://localhost:8080/api/v1/auth/google/callback" // In production, this should be an environment variable
)

// OAuth configuration for Google
var googleOAuthConfig = oauth2.Config{
	ClientID:     googleClientID,
	ClientSecret: googleClientSecret,
	RedirectURL:  googleRedirectURL,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

// LoginRequest defines the structure for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// GoogleAuthRequest defines the structure for Google authentication
type GoogleAuthRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// PasswordResetRequest defines the structure for password reset
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RefreshTokenRequest defines the structure for refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GoogleUserInfo defines the structure of Google user info
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Sub           string `json:"sub"` // Subject identifier
}

// TokenResponse defines the structure of the token response
type TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}

// Claims defines the structure of the JWT token
type Claims struct {
	UserID uint         `json:"user_id"`
	Email  string       `json:"email"`
	Roles  models.Roles `json:"roles"`
	Admin  bool         `json:"admin"`
	jwt.RegisteredClaims
}

// LoginUser handles user login with email and password
func LoginUser(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if result := database.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if user has a local password
	if !user.HasLocalPassword || user.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This account uses Google Sign-In. Please log in with Google."})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		// Create audit log for failed login
		createAuthAudit(c, user.ID, models.ActionLogin, false, "Invalid password")

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Update last login time
	user.LastLogin = time.Now()
	database.DB.Save(&user)

	// Create audit log for successful login
	createAuthAudit(c, user.ID, models.ActionLogin, true, "")

	// Generate tokens
	tokenResponse, err := generateTokens(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Return tokens and user info
	c.JSON(http.StatusOK, tokenResponse)
}

// AuthMiddleware verifies JWT token for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header has the format "Bearer <token>"
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := authHeader[7:]

		// Parse and validate token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("isAdmin", claims.Admin)

		c.Next()
	}
}
