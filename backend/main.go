package main

import (
	"log"
	"mis-system/database"
	"mis-system/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	database.ConnectDatabase()

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.RegisterUser)
			auth.POST("/login", handlers.LoginUser)
			auth.POST("/google", handlers.GoogleAuth)
			auth.GET("/google/login", handlers.GoogleLogin)
			auth.GET("/google/callback", handlers.GoogleCallback)
			auth.POST("/refresh", handlers.RefreshToken)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/forgot-password", handlers.ForgotPassword)
			auth.POST("/reset-password", handlers.ResetPassword)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(handlers.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/", handlers.GetAllUsers)
				users.GET("/:id", handlers.GetUserByID)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
			}

			// Me endpoint for getting current user info
			protected.GET("/me", handlers.GetCurrentUser)
		}
	}

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
