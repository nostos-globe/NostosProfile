package db

import (
	"main/internal/models"
	"main/internal/models/dto"

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

func (repo *FollowRepository) ListFollowers(profileID uint) ([]dto.ProfileBasicDTO, error) {
	var followers []dto.ProfileBasicDTO

	err := repo.DB.Table("auth.followers").
		Select("auth.profiles.profile_id, auth.profiles.username, auth.profiles.profile_picture").
		Joins("JOIN auth.profiles ON auth.followers.follower_id = auth.profiles.profile_id").
		Where("followed_id = ?", profileID).
		Scan(&followers).Error

	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (repo *FollowRepository) ListFollowing(profileID uint) ([]dto.ProfileBasicDTO, error) {
	var following []dto.ProfileBasicDTO

	err := repo.DB.Table("auth.followers").
		Select("auth.profiles.profile_id, auth.profiles.username, auth.profiles.profile_picture").
		Joins("JOIN auth.profiles ON auth.followers.followed_id = auth.profiles.profile_id").
		Where("follower_id = ?", profileID).
		Scan(&following).Error

	if err != nil {
		return nil, err
	}

	return following, nil
}
