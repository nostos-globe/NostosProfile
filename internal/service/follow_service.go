package service

import (
	"main/internal/db"
	"main/internal/models"
	"main/internal/models/dto"
)

type FollowService struct {
	FollowRepo *db.FollowRepository
}

type FollowersResponse struct {
	Count    int                   `json:"count"`
	Profiles []dto.ProfileBasicDTO `json:"profiles"`
}

func (s *FollowService) GetFollowByIDs(followerID uint, followedID uint) (any, error) {
	return s.FollowRepo.GetFollowByIDs(followerID, followedID)
}

func (s *FollowService) FollowUser(follow *models.Follow) error {
	return s.FollowRepo.FollowUser(follow)
}

func (s *FollowService) UnFollowUser(follow *models.Follow) error {
	return s.FollowRepo.UnFollowUser(follow)
}

func (s *FollowService) ListFollowers(profileID uint) (any, error) {
	followers, err := s.FollowRepo.ListFollowers(profileID)
	if err != nil {
		return nil, err
	}

	response := FollowersResponse{
		Count:    len(followers),
		Profiles: followers,
	}

	return response, nil
}

func (s *FollowService) ListFollowing(profileID uint) (any, error) {
	following, err := s.FollowRepo.ListFollowing(profileID)
	if err != nil {
		return nil, err
	}

	response := FollowersResponse{
		Count:    len(following),
		Profiles: following,
	}

	return response, nil
}
