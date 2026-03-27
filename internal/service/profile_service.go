package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository"
)

var ErrProfileAlreadyExists = errors.New("profile already exists")

type CreateProfileInput struct {
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type ProfileService interface {
	CreateProfile(ctx context.Context, userID int64, input CreateProfileInput) (*model.Profile, error)
}

type profileService struct {
	profileRepo repository.ProfileRepository
}

func NewProfileService(profileRepo repository.ProfileRepository) ProfileService {
	return &profileService{profileRepo: profileRepo}
}

func (s *profileService) CreateProfile(ctx context.Context, userID int64, input CreateProfileInput) (*model.Profile, error) {
	existing, err := s.profileRepo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrProfileAlreadyExists
	}

	profile := &model.Profile{
		UserID:    userID,
		Bio:       input.Bio,
		AvatarURL: input.AvatarURL,
	}

	if err = s.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("create profile: %w", err)
	}
	return profile, nil
}
