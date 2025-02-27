package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controller "main/internal/api"
	dbRepo "main/internal/db"
	"main/internal/service"
	"main/pkg/config"
	"main/pkg/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	secretsManager := config.GetSecretsManager()
	if secretsManager != nil {
		secrets := secretsManager.LoadSecrets()
		for key, value := range secrets {
			os.Setenv(key, value)
		}
	} else {
		log.Println("Falling back to environment variables")
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize repositories
	profileRepo := &dbRepo.ProfileRepository{DB: database}
	// Initialize authClient
	authClient := &service.AuthClient{BaseURL: cfg.AuthServiceUrl}
	// Initialize services
	profileService := &service.ProfileService{ProfileRepo: profileRepo}
	// Initialize controllers
	profileHandler := &controller.ProfileController{ProfileService: profileService, AuthClient: authClient}

	// Initialize Gin
	r := gin.Default()

	// Profile routes
	api := r.Group("/api")
	{
		// Protected profile routes
		api.POST("/profiles", profileHandler.CreateProfile)
		//api.PUT("/profiles", profileHandler.UpdateProfile)

		// Public profile routes
		api.GET("/profiles/user/:userID", profileHandler.GetProfile)
		//api.GET("/profiles/username/:username", profileHandler.GetProfileByUsername)
		//api.GET("/profiles/search", profileHandler.SearchProfiles)
	}

	// Start server
	log.Println("Server running on http://localhost:8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
