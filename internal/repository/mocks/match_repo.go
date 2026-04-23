package mocks

import (
	"context"

	"github.com/Fankemp/GameMatch/internal/model"
)

type MatchRepo struct {
	Matches    []*model.Match
	MatchExists bool // control exists response
	NextID     int64
}

func NewMatchRepo() *MatchRepo {
	return &MatchRepo{NextID: 1}
}

func (r *MatchRepo) Create(_ context.Context, match *model.Match) error {
	match.ID = r.NextID
	r.NextID++
	r.Matches = append(r.Matches, match)
	return nil
}

func (r *MatchRepo) GetByUserID(_ context.Context, userID int64) ([]*model.MatchWithUser, error) {
	return []*model.MatchWithUser{}, nil
}

func (r *MatchRepo) Exists(_ context.Context, userAID, userBID int64, gameID string) (bool, error) {
	return r.MatchExists, nil
}
