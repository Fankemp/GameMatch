package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Language string `json:"language"`
	Discord  string `json:"discord"`
	Telegram string `json:"telegram"`
	Region   string `json:"region"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*model.User, error)
	Login(ctx context.Context, input LoginInput) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*model.User, error) {
	existing, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
		Age:          input.Age,
		Language:     input.Language,
		Discord:      input.Discord,
		Telegram:     input.Telegram,
		Region:       input.Region,
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}
