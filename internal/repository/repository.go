package repository

import (
	"context"

	"github.com/Fankemp/GameMatch/internal/model"
)

type CardRepository interface {
	Create(ctx context.Context, card *model.GameCard) error
	GetByID(ctx context.Context, id int64) (*model.GameCard, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.GameCard, error)
	Update(ctx context.Context, card *model.GameCard) error
	Delete(ctx context.Context, id, userID int64) error
	GetFeed(ctx context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error)
}

type MatchRepository interface {
	Create(ctx context.Context, match *model.Match) error
	GetByUserID(ctx context.Context, userID int64) ([]*model.MatchWithUser, error)
	Exists(ctx context.Context, userAID, userBID int64, gameID string) (bool, error)
}

type ProfileRepository interface {
	Create(ctx context.Context, profile *model.Profile) error
	GetByUserID(ctx context.Context, userID int64) (*model.Profile, error)
	GetByID(ctx context.Context, id int64) (*model.Profile, error)
	Update(ctx context.Context, profile *model.Profile) error
}

type SwipeRepository interface {
	Create(ctx context.Context, swipe *model.Swipe) error
	Exists(ctx context.Context, userID, targetCardID int64) (bool, error)
	CheckMutualLike(ctx context.Context, targetUserID, currentUserID int64, gameID string) (bool, error)
}
