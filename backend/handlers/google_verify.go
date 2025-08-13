package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GoogleVerifyRequest represents the request to verify a Google ID token
type GoogleVerifyRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// VerifyGoogleToken verifies a Google ID token and returns user info without creating a session
func VerifyGoogleToken(c *gin.Context) {
	var req GoogleVerifyRequest
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

	// Return user info without creating a session
	userInfo := gin.H{
		"email":      tokenInfo["email"].(string),
		"sub":        tokenInfo["sub"].(string),
		"givenName":  tokenInfo["given_name"].(string),
		"familyName": tokenInfo["family_name"].(string),
		"name":       tokenInfo["name"].(string),
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userInfo,
	})
}
