package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository"
)

var (
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrProfileNotFound      = errors.New("profile not found")
)

type CreateProfileInput struct {
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateProfileInput struct {
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type ProfileService interface {
	CreateProfile(ctx context.Context, userID int64, input CreateProfileInput) (*model.Profile, error)
	UpdateProfile(ctx context.Context, userID int64, input UpdateProfileInput) (*model.Profile, error)
	GetMyProfile(ctx context.Context, userID int64) (*model.Profile, error)
	GetProfileByID(ctx context.Context, id int64) (*model.Profile, error)
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

func (s *profileService) UpdateProfile(ctx context.Context, userID int64, input UpdateProfileInput) (*model.Profile, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrProfileNotFound
	}

	profile.Bio = input.Bio
	profile.AvatarURL = input.AvatarURL

	if err = s.profileRepo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return profile, nil
}

func (s *profileService) GetMyProfile(ctx context.Context, userID int64) (*model.Profile, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}

func (s *profileService) GetProfileByID(ctx context.Context, id int64) (*model.Profile, error) {
	profile, err := s.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}
