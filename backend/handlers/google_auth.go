package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mis-system/database"
	"mis-system/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GoogleLogin initiates the Google OAuth flow
func GoogleLogin(c *gin.Context) {
	// Generate a random state to prevent CSRF
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	// Store state in cookie for verification later
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	// Redirect to Google's OAuth page
	url := googleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth callback from Google
func GoogleCallback(c *gin.Context) {
	// Verify state to prevent CSRF
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Exchange authorization code for token
	code := c.Query("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Get user info from Google
	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info from Google"})
		return
	}
	defer resp.Body.Close()

	googleUser := GoogleUserInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode Google user info"})
		return
	}

	// Process Google user info
	tokenResponse, err := processGoogleUser(c, &googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return tokens and user info
	c.JSON(http.StatusOK, tokenResponse)
}

// GoogleAuth handles direct Google authentication with ID token
func GoogleAuth(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify ID token with Google
	tokenInfo, err := verifyGoogleIDToken(req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google ID token"})
		return
	}

	// Convert tokenInfo to GoogleUserInfo
	googleUser := GoogleUserInfo{
		Sub:        tokenInfo["sub"].(string),
		Email:      tokenInfo["email"].(string),
		Name:       tokenInfo["name"].(string),
		GivenName:  tokenInfo["given_name"].(string),
		FamilyName: tokenInfo["family_name"].(string),
	}

	if verifiedEmail, ok := tokenInfo["email_verified"].(bool); ok {
		googleUser.VerifiedEmail = verifiedEmail
	}

	if picture, ok := tokenInfo["picture"].(string); ok {
		googleUser.Picture = picture
	}

	// Process Google user info
	tokenResponse, err := processGoogleUser(c, &googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return tokens and user info
	c.JSON(http.StatusOK, tokenResponse)
}

// processGoogleUser handles the common processing for Google users
func processGoogleUser(c *gin.Context, googleUser *GoogleUserInfo) (*TokenResponse, error) {
	// Look for existing user by Google Sub ID
	var user models.User
	result := database.DB.Where("google_sub = ?", googleUser.Sub).First(&user)

	// Check if user exists by email if not found by Google ID
	if result.Error != nil {
		result = database.DB.Where("email = ?", googleUser.Email).First(&user)
	}

	// Create new user if not found
	if result.Error != nil {
		// Auto-provision a new user linked to Google
		names := strings.Fields(googleUser.Name)
		firstName := googleUser.GivenName
		lastName := googleUser.FamilyName

		// Fallback to parsed name if given/family name not provided
		if firstName == "" && len(names) > 0 {
			firstName = names[0]
		}
		if lastName == "" && len(names) > 1 {
			lastName = strings.Join(names[1:], " ")
		}

		// Create new user
		user = models.User{
			Email:            googleUser.Email,
			GoogleSub:        googleUser.Sub,
			FirstName:        firstName,
			LastName:         lastName,
			HasLocalPassword: false,
			IsActive:         true,
			IsAdmin:          false,
			Roles:            models.Roles{models.RoleUser}, // Default role
			LastLogin:        time.Now(),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			return nil, err
		}

		// Create audit log for new user
		createAuthAudit(c, user.ID, models.ActionRegister, true, "Google account auto-provisioned")
	} else {
		// Update existing user with Google info
		user.GoogleSub = googleUser.Sub
		user.LastLogin = time.Now()

		// Only update name if previously empty
		if user.FirstName == "" && googleUser.GivenName != "" {
			user.FirstName = googleUser.GivenName
		}
		if user.LastName == "" && googleUser.FamilyName != "" {
			user.LastName = googleUser.FamilyName
		}

		if err := database.DB.Save(&user).Error; err != nil {
			return nil, err
		}
	}

	// Create auth audit log
	createAuthAudit(c, user.ID, models.ActionGoogleAuth, true, "")

	// Generate tokens
	return generateTokens(c, &user)
}

// verifyGoogleIDToken verifies the Google ID token
func verifyGoogleIDToken(idToken string) (map[string]interface{}, error) {
	// Make a request to Google's token info endpoint
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to verify token: %s", resp.Status)
	}

	var tokenInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, err
	}

	// Verify audience (your client ID)
	if aud, ok := tokenInfo["aud"].(string); !ok || aud != googleClientID {
		return nil, fmt.Errorf("invalid token audience")
	}

	return tokenInfo, nil
}

// generateTokens creates and returns access and refresh tokens
func generateTokens(c *gin.Context, user *models.User) (*TokenResponse, error) {
	// Create access token
	accessTokenExp := time.Now().Add(accessTokenExp)
	accessTokenClaims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  user.Roles,
		Admin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	// Generate a secure random refresh token
	refreshTokenBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, refreshTokenBytes); err != nil {
		return nil, err
	}
	refreshTokenString := hex.EncodeToString(refreshTokenBytes)

	// Hash the refresh token for storage
	hasher := sha256.New()
	hasher.Write([]byte(refreshTokenString))
	hashedRefreshToken := hex.EncodeToString(hasher.Sum(nil))

	// Get client info
	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		deviceID = uuid.New().String()
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	// Create a new session record
	session := models.Session{
		UserID:       user.ID,
		RefreshToken: hashedRefreshToken,
		DeviceID:     deviceID,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(refreshTokenExp),
	}

	if err := database.DB.Create(&session).Error; err != nil {
		return nil, err
	}

	// Return token response
	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		User: gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"roles":     user.Roles,
			"isAdmin":   user.IsAdmin,
		},
	}, nil
}

// RefreshToken handles refresh token requests
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the refresh token to compare with stored hash
	hasher := sha256.New()
	hasher.Write([]byte(req.RefreshToken))
	hashedRefreshToken := hex.EncodeToString(hasher.Sum(nil))

	// Find the session
	var session models.Session
	if err := database.DB.Where("refresh_token = ? AND expires_at > ? AND revoked_at IS NULL",
		hashedRefreshToken, time.Now()).First(&session).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Get the user
	var user models.User
	if err := database.DB.First(&user, session.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Revoke the old refresh token
	session.RevokedAt = time.Now()
	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke old token"})
		return
	}

	// Create auth audit log
	createAuthAudit(c, user.ID, models.ActionRefresh, true, "")

	// Generate new tokens
	tokenResponse, err := generateTokens(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// Logout handles user logout
func Logout(c *gin.Context) {
	// Get refresh token from request
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the refresh token
	hasher := sha256.New()
	hasher.Write([]byte(req.RefreshToken))
	hashedRefreshToken := hex.EncodeToString(hasher.Sum(nil))

	// Find and revoke the session
	var session models.Session
	if err := database.DB.Where("refresh_token = ? AND revoked_at IS NULL",
		hashedRefreshToken).First(&session).Error; err != nil {
		// Just return success even if token not found for security
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		return
	}

	// Revoke the session
	session.RevokedAt = time.Now()
	database.DB.Save(&session)

	// Create audit log
	createAuthAudit(c, session.UserID, models.ActionLogout, true, "")

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ForgotPassword initiates password reset
func ForgotPassword(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal whether the email exists or not
		c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you'll receive password reset instructions"})
		return
	}

	// Check if user has a local password
	if !user.HasLocalPassword {
		// For security, don't reveal that the user doesn't have a local password
		c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you'll receive password reset instructions"})
		return
	}

	// Generate a password reset token (in a real app, this would be stored and used to verify the reset request)
	// For now we just generate it but don't use it since we're not implementing the full email flow
	_ = uuid.New().String()

	// In a real application, send an email with a link including this token
	// For now, just imagine we've sent an email

	// Create audit log
	createAuthAudit(c, user.ID, models.ActionPasswordReset, true, "Password reset requested")

	c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you'll receive password reset instructions"})
}

// createAuthAudit creates an auth audit log entry
func createAuthAudit(c *gin.Context, userID uint, action models.AuditAction, success bool, details string) {
	audit := models.AuthAudit{
		UserID:    userID,
		Action:    action,
		Success:   success,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		DeviceID:  c.GetHeader("X-Device-ID"),
		Details:   details,
	}

	database.DB.Create(&audit)
}
