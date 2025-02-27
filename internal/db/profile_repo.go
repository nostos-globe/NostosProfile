package db

import (
	"main/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

// Basic CRUD operations for Profile
func (repo *ProfileRepository) CreateProfile(profile *models.Profile) error {
	return repo.DB.Table("auth.profiles").Create(profile).Error
}

func (repo *ProfileRepository) GetProfileByUserID(userID uint) (*models.Profile, error) {
	var profile models.Profile
	err := repo.DB.Table("auth.profiles").Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (repo *ProfileRepository) GetProfileByUsername(username string) (*models.Profile, error) {
	var profile models.Profile
	err := repo.DB.Table("auth.profiles").Where("username = ?", username).First(&profile).Error
	return &profile, err
}

func (repo *ProfileRepository) UpdateProfile(profile *models.Profile) error {
	return repo.DB.Table("auth.profiles").Save(profile).Error
}

// Search profiles
func (repo *ProfileRepository) SearchProfiles(query string) ([]models.Profile, error) {
	var profiles []models.Profile
	err := repo.DB.Table("auth.profiles").Where("username ILIKE ?", "%"+query+"%").Find(&profiles).Error
	return profiles, err
}
