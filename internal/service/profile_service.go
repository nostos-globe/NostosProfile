package service

import (
	"fmt"
	"log"
	"main/internal/db"
	"main/internal/events"
	"main/internal/models"
	"time"
	"mime/multipart"
)

type ProfileService struct {
	ProfileRepo *db.ProfileRepository
	MinioClient *MinioService
	NatsClient  *events.NatsClient
	AuthClient  *AuthClient
}

// Update CreateProfile method
func (s *ProfileService) CreateProfile(profile *models.Profile, email string) error {
	err := s.ProfileRepo.CreateProfile(profile)
	if err != nil {
		return err
	}
	// Publish event after successful profile creation
	if s.NatsClient != nil {
		err = s.NatsClient.PublishUserRegistered(profile.UserID, email, profile.Username)
		if err != nil {
			log.Printf("Warning: Failed to publish user.registered event: %v", err)
			// Don't return error as this is not critical for profile creation
		}
	}
	return nil
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

func (s *ProfileService) UploadAvatar(file multipart.File, header *multipart.FileHeader) (string, error) {
    return s.MinioClient.UploadFile(file, header)
}

func (s *ProfileService) GetAvatarURL(filename string) (string, error) {
    duration := time.Hour * 1
    return s.MinioClient.GetPresignedURL(filename, duration)
}
