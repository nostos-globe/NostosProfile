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
	// Parse multipart form
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	var profile models.Profile
	
	// Get form values
    profile.Username = ctx.Request.FormValue("username")
    bio := ctx.Request.FormValue("bio")
    profile.Bio = &bio
	if theme := ctx.Request.FormValue("theme"); theme != "" {
		profile.Theme = &theme
	}
	if website := ctx.Request.FormValue("website"); website != "" {
		profile.Website = &website
	}
	if language := ctx.Request.FormValue("language"); language != "" {
		profile.Language = &language
	}
	if birthdate := ctx.Request.FormValue("birthdate"); birthdate != "" {
		profile.Birthdate = &birthdate
	}

	// Handle file upload
	file, header, err := ctx.Request.FormFile("profilePicture")
	if err == nil && file != nil {
		defer file.Close()
		
		filename, err := c.ProfileService.UploadAvatar(file, header)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload avatar"})
			return
		}
		profile.ProfilePicture = &filename
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

	email, err := c.AuthClient.GetUserEmail(tokenCookie)
	if err != nil || TokenResponse == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	profile.UserID = TokenResponse
	if err := c.ProfileService.CreateProfile(&profile, email); err != nil {
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
    if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
        return
    }

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

	// Update profile fields from form values
    if username := ctx.Request.FormValue("username"); username != "" {
        existingProfile.Username = username
    }
    if bio := ctx.Request.FormValue("bio"); bio != "" {
        existingProfile.Bio = &bio
    }
	
	if theme := ctx.Request.FormValue("theme"); theme != "" {
		existingProfile.Theme = &theme
	}
	if website := ctx.Request.FormValue("website"); website != "" {
		existingProfile.Website = &website
	}
	if language := ctx.Request.FormValue("language"); language != "" {
		existingProfile.Language = &language
	}
	if birthdate := ctx.Request.FormValue("birthdate"); birthdate != "" {
		existingProfile.Birthdate = &birthdate
	}

	// Handle new profile picture
	file, header, err := ctx.Request.FormFile("profilePicture")
    if err == nil && file != nil {
        defer file.Close()
        
        // Delete old profile picture if exists
        if existingProfile.ProfilePicture != nil {
            _ = c.ProfileService.MinioClient.DeleteObject(*existingProfile.ProfilePicture)  // Changed MinioService to MinioClient
        }
        
        filename, err := c.ProfileService.UploadAvatar(file, header)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload avatar"})
            return
        }
        existingProfile.ProfilePicture = &filename
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

	if err := c.ProfileService.DeleteProfile(existingProfile); err != nil {
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

	if profile != nil && profile.ProfilePicture != nil {
		url, err := c.ProfileService.GetAvatarURL(*profile.ProfilePicture)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate avatar URL"})
			return
		}
		profile.ProfilePicture = &url
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

	for i, profile := range profiles {
		if profile.ProfilePicture != nil {
			url, err := c.ProfileService.GetAvatarURL(*profile.ProfilePicture)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate avatar URL"})
				return
			}
			profiles[i].ProfilePicture = &url
		}
	}

	ctx.JSON(http.StatusOK, profiles)
}

func (c *ProfileController) GetProfileAvatar(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	// Validate token
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

	// Get profile to verify the avatar exists
	profile, err := c.ProfileService.GetProfileByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil || profile.ProfilePicture == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "avatar not found"})
		return
	}

	// Get presigned URL from MinIO using the stored filename
	url, err := c.ProfileService.GetAvatarURL(*profile.ProfilePicture)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate avatar URL"})
		return
	}

	// Return the presigned URL
	ctx.JSON(http.StatusOK, gin.H{"url": url})
}