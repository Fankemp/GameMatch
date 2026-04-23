package mocks

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
)

type CardRepo struct {
	Cards  map[int64]*model.GameCard
	NextID int64
}

func NewCardRepo() *CardRepo {
	return &CardRepo{Cards: make(map[int64]*model.GameCard), NextID: 1}
}

func (r *CardRepo) Create(_ context.Context, card *model.GameCard) error {
	card.ID = r.NextID
	card.IsActive = true
	r.NextID++
	r.Cards[card.ID] = card
	return nil
}

func (r *CardRepo) GetByID(_ context.Context, id int64) (*model.GameCard, error) {
	c, ok := r.Cards[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return c, nil
}

func (r *CardRepo) GetByUserID(_ context.Context, userID int64) ([]*model.GameCard, error) {
	var result []*model.GameCard
	for _, c := range r.Cards {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *CardRepo) Update(_ context.Context, card *model.GameCard) error {
	if _, ok := r.Cards[card.ID]; !ok {
		return fmt.Errorf("not found")
	}
	r.Cards[card.ID] = card
	return nil
}

func (r *CardRepo) Delete(_ context.Context, id, userID int64) error {
	c, ok := r.Cards[id]
	if !ok || c.UserID != userID {
		return fmt.Errorf("card not found")
	}
	delete(r.Cards, id)
	return nil
}

func (r *CardRepo) GetFeed(_ context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error) {
	return []*model.FeedCard{}, nil
}
