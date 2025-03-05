package service

import (
	"main/internal/db"
	"main/internal/models"
)

type FollowService struct {
	FollowRepo *db.FollowRepository
}

func (s *FollowService) CreateFollow(follow *models.Follow) error {
	return s.FollowRepo.CreateFollow(follow)
}
