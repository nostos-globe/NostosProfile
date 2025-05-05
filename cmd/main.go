package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controller "main/internal/api"
	dbRepo "main/internal/db"
	"main/internal/events"
	"main/internal/service"
	"main/pkg/config"
	"main/pkg/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	minioManager := config.InitMinIO()
	if minioManager == nil {
		log.Println("Faliing to init MinIO")
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
	followRepo := &dbRepo.FollowRepository{DB: database}
	// Initialize authClient
	authClient := &service.AuthClient{BaseURL: cfg.AuthServiceUrl}
	// Initialize NATS client
	natsClient, err := events.NewNatsClient(cfg.NatsURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS: %v", err)
		// Continue without NATS if connection fails
	}

	// Initialize services
	minioService := service.NewMinioService()
	profileService := &service.ProfileService{ProfileRepo: profileRepo, MinioClient: minioService, NatsClient: natsClient, AuthClient: authClient}
	followService := &service.FollowService{
		FollowRepo:     followRepo,
		ProfileService: profileService,
	}
	// Initialize controllers
	profileHandler := &controller.ProfileController{ProfileService: profileService, AuthClient: authClient}
	followHandler := &controller.FollowController{FollowService: followService, ProfileService: profileService, AuthClient: authClient}

	// Initialize Gin
	r := gin.Default()

	// Profile routes
	api := r.Group("/api")
	{
		// Protected profile routes
		api.POST("/profiles", profileHandler.CreateProfile)
		api.POST("/profiles/update", profileHandler.UpdateProfile)
		api.POST("/profiles/updateProfileByID", profileHandler.UpdateProfileByID)
		api.POST("/profiles/delete", profileHandler.DeleteProfile)

		// Public profile routes
		api.GET("/profiles/user/:userID", profileHandler.GetProfile)
		api.GET("/profiles/username/:username", profileHandler.GetProfileByUsername)
		api.POST("/profiles/search", profileHandler.SearchProfiles)
		api.GET("/profiles/userAvatar/:userID", profileHandler.GetProfileAvatar)

		// Followers
		api.POST("/follow/:followedID", followHandler.FollowUser)
		api.POST("/unfollow/:followedID", followHandler.UnFollowUser)
		api.GET("/:profileID/followers", followHandler.ListFollowers)
		api.GET("/:profileID/following", followHandler.ListFollowing)

	}

	// Start server
	log.Println("Server running on http://localhost:8082")
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
