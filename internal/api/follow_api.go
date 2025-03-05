package controller

import (
	"main/internal/models"
	"main/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	ProfileService 	*service.ProfileService
	FollowService 	*service.FollowService
	AuthClient     	*service.AuthClient
}

func (c *FollowController) FollowUser(ctx *gin.Context) {
    // Get user ID from token
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

    // Get follower profile using user ID
    followerProfile, err := c.ProfileService.GetProfileByUserID(userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get follower profile"})
        return
    }

    // Get followed user ID from params
    followedID, err := strconv.Atoi(ctx.Param("followed_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid followed_id parameter"})
        return
    }

    // Get followed profile using followed user ID
    followedProfile, err := c.ProfileService.GetProfileByUserID(uint(followedID))
    if err!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followed profile"})
        return
    }
	
    // Create follow relationship
	var follow models.Follow

	follow.FollowerID = followerProfile.ID
	follow.FollowedID = followedProfile.ID

	if err := c.FollowService.CreateFollow(&follow); err!= nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}
