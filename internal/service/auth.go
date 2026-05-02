package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
)

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type authService struct {
	userRepo UserRepository
	hasher   PasswordHasher
	tokens   TokenManager
}

func NewAuthService(userRepo UserRepository, h PasswordHasher, tm TokenManager) AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   h,
		tokens:   tm,
	}
}

func (s *authService) SignUp(ctx context.Context, input SignUpInput) (*AuthResponse, error) {
	existing, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}

	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
		Age:          input.Age,
		Language:     input.Language,
		Discord:      input.Discord,
		Telegram:     input.Telegram,
		Region:       input.Region,
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := s.tokens.NewJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user}, nil
}

func (s *authService) SignIn(ctx context.Context, input SignInInput) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.hasher.Compare(user.PasswordHash, input.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokens.NewJWT(user.ID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &AuthResponse{Token: token, User: user}, nil
}

func (s *authService) GetMe(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
