package db

import (
	"main/internal/models"

	"gorm.io/gorm"
)

type FollowRepository struct {
	DB *gorm.DB
}

func (repo *FollowRepository) CreateFollow(follow *models.Follow) error {
	return repo.DB.Table("auth.followers").Create(follow).Error
}
