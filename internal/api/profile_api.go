package controller

import (
	"main/internal/models"
	"main/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	ProfileService *service.ProfileService
	AuthClient     *service.AuthClient
}

func (c *ProfileController) CreateProfile(ctx *gin.Context) {
	var profile models.Profile
	if err := ctx.ShouldBindJSON(&profile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from authenticated context
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || TokenResponse == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	profile.UserID = TokenResponse
	if err := c.ProfileService.CreateProfile(&profile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, profile)
}

func (c *ProfileController) GetProfileByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.ValidateToken(tokenCookie)
	if err != nil || TokenResponse == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to validate the token"})
		return
	}

	profile, err := c.ProfileService.GetProfileByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	var profile models.Profile
	if err := ctx.ShouldBindJSON(&profile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get user ID from token"})
		return
	}

	existingProfile, err := c.ProfileService.GetProfileByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingProfile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	if profile.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username cannot be empty"})
		return
	}

	existingProfile.Username = profile.Username
	existingProfile.Bio = profile.Bio

	if profile.ProfilePicture != nil {
		existingProfile.ProfilePicture = profile.ProfilePicture
	}

	if profile.Theme != nil {
		existingProfile.Theme = profile.Theme
	}

	if profile.Birthdate != nil {
		existingProfile.Birthdate = profile.Birthdate
	}

	if profile.Language != nil {
		existingProfile.Language = profile.Language
	}

	if profile.Website != nil {
		existingProfile.Website = profile.Website
	}

	if err := c.ProfileService.UpdateProfile(existingProfile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, existingProfile)
}

func (c *ProfileController) UpdateProfileByID(ctx *gin.Context) {
	var profile models.Profile
	if err := ctx.ShouldBindJSON(&profile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.ValidateToken(tokenCookie)
	if err != nil || TokenResponse == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to validate the token"})
		return
	}

	existingProfile, err := c.ProfileService.GetProfileByUserID(profile.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingProfile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	if profile.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username cannot be empty"})
		return
	}
	existingProfile.Username = profile.Username
	existingProfile.Bio = profile.Bio

	if profile.ProfilePicture != nil {
		existingProfile.ProfilePicture = profile.ProfilePicture
	}

	if profile.Theme != nil {
		existingProfile.Theme = profile.Theme
	}

	if profile.Birthdate != nil {
		existingProfile.Birthdate = profile.Birthdate
	}

	if profile.Language != nil {
		existingProfile.Language = profile.Language
	}

	if profile.Website != nil {
		existingProfile.Website = profile.Website
	}

	if err := c.ProfileService.UpdateProfile(existingProfile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, existingProfile)
}

func (c *ProfileController) DeleteProfile(ctx *gin.Context) {
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get user ID from token"})
		return
	}

	existingProfile, err := c.ProfileService.GetProfileByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingProfile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	if err := c.ProfileService.DeleteProfile(existingProfile); err!= nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

func (c *ProfileController) GetProfile(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.ValidateToken(tokenCookie)
	if err != nil || TokenResponse == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to validate the token"})
		return
	}

	// Convert string to uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	profile, err := c.ProfileService.GetProfileByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func (c *ProfileController) SearchProfiles(ctx *gin.Context) {
	var searchRequest struct {
		Query string `json:"query"`
	}
	if err := ctx.ShouldBindJSON(&searchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.ValidateToken(tokenCookie)
	if err != nil || TokenResponse == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to validate the token"})
		return
	}

	profiles, err := c.ProfileService.SearchProfiles(searchRequest.Query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profiles)
}
