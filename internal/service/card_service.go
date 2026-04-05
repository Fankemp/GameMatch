package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository"
)

var ErrCardNotFound = errors.New("card not found")

type CreateCardInput struct {
	GameID      string `json:"game_id"`
	Rank        string `json:"rank"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

type UpdateCardInput struct {
	Rank        string `json:"rank"`
	Role        string `json:"role"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type CardService interface {
	Create(ctx context.Context, userID int64, input CreateCardInput) (*model.GameCard, error)
	GetMyCards(ctx context.Context, userID int64) ([]*model.GameCard, error)
	Update(ctx context.Context, userID, cardID int64, input UpdateCardInput) (*model.GameCard, error)
	Delete(ctx context.Context, userID, cardID int64) error
}

type cardService struct {
	cardRepo repository.CardRepository
}

func NewCardService(cardRepo repository.CardRepository) CardService {
	return &cardService{cardRepo: cardRepo}
}

func (s *cardService) Create(ctx context.Context, userID int64, input CreateCardInput) (*model.GameCard, error) {
	card := &model.GameCard{
		UserID:      userID,
		GameID:      input.GameID,
		Rank:        input.Rank,
		Role:        input.Role,
		Description: input.Description,
	}
	if err := s.cardRepo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}
	return card, nil
}

func (s *cardService) GetMyCards(ctx context.Context, userID int64) ([]*model.GameCard, error) {
	return s.cardRepo.GetByUserID(ctx, userID)
}

func (s *cardService) Update(ctx context.Context, userID, cardID int64, input UpdateCardInput) (*model.GameCard, error) {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}
	if card.UserID != userID {
		return nil, ErrCardNotFound
	}

	if input.Rank != "" {
		card.Rank = input.Rank
	}
	if input.Role != "" {
		card.Role = input.Role
	}
	if input.Description != "" {
		card.Description = input.Description
	}
	if input.IsActive != nil {
		card.IsActive = *input.IsActive
	}

	if err := s.cardRepo.Update(ctx, card); err != nil {
		return nil, fmt.Errorf("update card: %w", err)
	}
	return card, nil
}

func (s *cardService) Delete(ctx context.Context, userID, cardID int64) error {
	return s.cardRepo.Delete(ctx, cardID, userID)
}
