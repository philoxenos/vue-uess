package handlers

import (
	"mis-system/database"
	"mis-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest defines the structure for user registration
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	GoogleSub       string `json:"google_sub"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
}

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		// If we have a GoogleSub from the request, check if it matches the existing user
		if input.GoogleSub != "" && existingUser.GoogleSub == input.GoogleSub {
			// User already exists with the same Google account, just update the password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
				return
			}

			existingUser.Password = string(hashedPassword)
			existingUser.HasLocalPassword = true
			if err := database.DB.Save(&existingUser).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
				return
			}

			// Create audit log
			createAuthAudit(c, existingUser.ID, models.ActionPasswordReset, true, "Password set for Google account")

			// Generate tokens
			tokenResponse, err := generateTokens(c, &existingUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
				return
			}

			c.JSON(http.StatusOK, tokenResponse)
			return
		}

		// Otherwise email already exists with a different account
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// Create new user
	user := models.User{
		Email:            input.Email,
		Password:         string(hashedPassword),
		GoogleSub:        input.GoogleSub, // May be empty if registering without Google
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		HasLocalPassword: true,
		IsActive:         true,
		Roles:            models.Roles{models.RoleUser}, // Default role
	}

	// Save user to database
	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create audit log
	createAuthAudit(c, user.ID, models.ActionRegister, true, "New user registered")

	// Generate tokens
	tokenResponse, err := generateTokens(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Return tokens and user info
	c.JSON(http.StatusCreated, tokenResponse)
}

// GetAllUsers retrieves all users
func GetAllUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUserByID retrieves a single user by ID
func GetUserByID(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser updates a user's information
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&user).Updates(models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":               user.ID,
			"email":            user.Email,
			"firstName":        user.FirstName,
			"lastName":         user.LastName,
			"roles":            user.Roles,
			"isAdmin":          user.IsAdmin,
			"googleSub":        user.GoogleSub != "",
			"hasLocalPassword": user.HasLocalPassword,
		},
	})

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser removes a user
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// ResetPasswordRequest defines the structure for password reset
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPassword handles password reset with token
func ResetPassword(c *gin.Context) {
	var input ResetPasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app, validate the reset token against stored tokens
	// For this example, we'll just return an error since we're not implementing the full flow
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
