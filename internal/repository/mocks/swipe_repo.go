package mocks

import (
	"context"

	"github.com/Fankemp/GameMatch/internal/model"
)

type SwipeRepo struct {
	Swipes     []*model.Swipe
	MutualLike bool // control mutual like response
	NextID     int64
}

func NewSwipeRepo() *SwipeRepo {
	return &SwipeRepo{NextID: 1}
}

func (r *SwipeRepo) Create(_ context.Context, swipe *model.Swipe) error {
	swipe.ID = r.NextID
	r.NextID++
	r.Swipes = append(r.Swipes, swipe)
	return nil
}

func (r *SwipeRepo) Exists(_ context.Context, userID, targetCardID int64) (bool, error) {
	for _, s := range r.Swipes {
		if s.UserID == userID && s.TargetCardID == targetCardID {
			return true, nil
		}
	}
	return false, nil
}

func (r *SwipeRepo) CheckMutualLike(_ context.Context, targetUserID, currentUserID int64, gameID string) (bool, error) {
	return r.MutualLike, nil
}
