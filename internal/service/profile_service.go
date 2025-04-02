package service

import (
	"fmt"
	"main/internal/db"
	"main/internal/models"
	"time"
)

type ProfileService struct {
	ProfileRepo *db.ProfileRepository
	MinioClient *MinioService
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

func (s *ProfileService) UpdateProfile(profile *models.Profile) error {
	existingProfile, err := s.ProfileRepo.GetProfileByUserID(profile.UserID)
	if err != nil {
		return err
	}

	if existingProfile == nil {
		return fmt.Errorf("profile not found")
	}

	existingProfile = profile

	return s.ProfileRepo.UpdateProfile(existingProfile)
}

func (s *ProfileService) DeleteProfile(profile *models.Profile) error {
	return s.ProfileRepo.DeleteProfile(profile)
}
func (s *ProfileService) GetProfileByUsername(username string) (*models.Profile, error) {
	return s.ProfileRepo.GetProfileByUsername(username)
}

func (s *ProfileService) GetProfileByUserID(userID uint) (*models.Profile, error) {
	return s.ProfileRepo.GetProfileByUserID(userID)
}

func (s *ProfileService) SearchProfiles(query string) ([]*models.Profile, error) {
	profiles, err := s.ProfileRepo.SearchProfiles(query)
	if err != nil {
		return nil, err
	}

	// Convert []models.Profile to []*models.Profile
	profilePtrs := make([]*models.Profile, len(profiles))
	for i := range profiles {
		profilePtrs[i] = &profiles[i]
	}

	return profilePtrs, nil
}

func (s *ProfileService) GetAvatarURL(filename string) (string, error) {
    return s.MinioClient.GetPresignedURL(filename, time.Minute*5)
}