package service

import (
	"context"
	"errors"

	"github.com/Fankemp/GameMatch/internal/model"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternal           = errors.New("internal server error")
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type TokenManager interface {
	NewJWT(userID int64) (string, error)
	Parse(token string) (int64, error)
}

type AuthService interface {
	SignUp(ctx context.Context, input SignUpInput) (*model.User, error)
	SignIn(ctx context.Context, input SignInInput) (string, error)
	GetMe(ctx context.Context, userID int64) (*model.User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type ProfileService interface {
	CreateProfile(ctx context.Context, userID int64, input CreateProfileInput) (*model.Profile, error)
	UpdateProfile(ctx context.Context, userID int64, input UpdateProfileInput) (*model.Profile, error)
	GetMyProfile(ctx context.Context, userID int64) (*model.Profile, error)
	GetProfileByID(ctx context.Context, id int64) (*model.Profile, error)
}

type CardService interface {
	Create(ctx context.Context, userID int64, input CreateCardInput) (*model.GameCard, error)
	GetMyCards(ctx context.Context, userID int64) ([]*model.GameCard, error)
	Update(ctx context.Context, userID, cardID int64, input UpdateCardInput) (*model.GameCard, error)
	Delete(ctx context.Context, userID, cardID int64) error
}

type FeedService interface {
	GetFeed(ctx context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error)
}

type SwipeService interface {
	Swipe(ctx context.Context, userID int64, input SwipeInput) (*SwipeResult, error)
}
