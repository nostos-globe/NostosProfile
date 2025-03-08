package db

import (
	"main/internal/models"

	"gorm.io/gorm"
)

type FollowRepository struct {
	DB *gorm.DB
}

func (repo *FollowRepository) GetFollowByIDs(followerID uint, followedID uint) (any, error) {
	var follow models.Follow
	err := repo.DB.Table("auth.followers").Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

func (repo *FollowRepository) FollowUser(follow *models.Follow) error {
	return repo.DB.Table("auth.followers").Create(follow).Error
}

func (repo *FollowRepository) UnFollowUser(follow *models.Follow) error {
	return repo.DB.Table("auth.followers").Delete(follow).Error
}
