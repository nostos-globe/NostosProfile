package service

import (
	"fmt"
	"main/internal/db"
	"main/internal/models"
)

type ProfileService struct {
	ProfileRepo *db.ProfileRepository
}

func (s *ProfileService) CreateProfile(profile *models.Profile) error {

	if profile.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	existingUserProfile, err := s.ProfileRepo.GetProfileByUserID(profile.UserID)
	if err == nil && existingUserProfile != nil {
		return fmt.Errorf("user already has a profile")
	}

	existingProfile, err := s.ProfileRepo.GetProfileByUsername(profile.Username)
	if err == nil && existingProfile != nil {
		return fmt.Errorf("username already exists")
	}

	return s.ProfileRepo.CreateProfile(profile)
}

func (s *ProfileService) GetProfileByUserID(userID uint) (*models.Profile, error) {
	return s.ProfileRepo.GetProfileByUserID(userID)
}
