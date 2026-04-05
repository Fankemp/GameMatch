package service

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository"
)

type SwipeInput struct {
	TargetCardID int64  `json:"target_card_id"`
	Action       string `json:"action"`
}

type SwipeResult struct {
	Matched bool         `json:"matched"`
	Match   *model.Match `json:"match,omitempty"`
}

type SwipeService interface {
	Swipe(ctx context.Context, userID int64, input SwipeInput) (*SwipeResult, error)
}

type swipeService struct {
	swipeRepo repository.SwipeRepository
	cardRepo  repository.CardRepository
	matchRepo repository.MatchRepository
}

func NewSwipeService(
	swipeRepo repository.SwipeRepository,
	cardRepo repository.CardRepository,
	matchRepo repository.MatchRepository,
) SwipeService {
	return &swipeService{
		swipeRepo: swipeRepo,
		cardRepo:  cardRepo,
		matchRepo: matchRepo,
	}
}

func (s *swipeService) Swipe(ctx context.Context, userID int64, input SwipeInput) (*SwipeResult, error) {
	// Check if already swiped
	exists, err := s.swipeRepo.Exists(ctx, userID, input.TargetCardID)
	if err != nil {
		return nil, fmt.Errorf("check swipe: %w", err)
	}
	if exists {
		return &SwipeResult{Matched: false}, nil
	}

	// Get target card to know the owner and game
	targetCard, err := s.cardRepo.GetByID(ctx, input.TargetCardID)
	if err != nil {
		return nil, fmt.Errorf("get target card: %w", err)
	}

	// Can't swipe on own card
	if targetCard.UserID == userID {
		return nil, fmt.Errorf("cannot swipe on own card")
	}

	// Save the swipe
	swipe := &model.Swipe{
		UserID:       userID,
		TargetCardID: input.TargetCardID,
		Action:       input.Action,
	}
	if err := s.swipeRepo.Create(ctx, swipe); err != nil {
		return nil, fmt.Errorf("create swipe: %w", err)
	}

	// If dislike, no match possible
	if input.Action != model.SwipeActionLike {
		return &SwipeResult{Matched: false}, nil
	}

	// Check if target user has also liked one of current user's cards for same game
	mutual, err := s.swipeRepo.CheckMutualLike(ctx, targetCard.UserID, userID, targetCard.GameID)
	if err != nil {
		return nil, fmt.Errorf("check mutual: %w", err)
	}

	if !mutual {
		return &SwipeResult{Matched: false}, nil
	}

	// Check if match already exists
	matchExists, err := s.matchRepo.Exists(ctx, userID, targetCard.UserID, targetCard.GameID)
	if err != nil {
		return nil, fmt.Errorf("check match exists: %w", err)
	}
	if matchExists {
		return &SwipeResult{Matched: false}, nil
	}

	// Create match
	match := &model.Match{
		UserAID: userID,
		UserBID: targetCard.UserID,
		GameID:  targetCard.GameID,
	}
	if err := s.matchRepo.Create(ctx, match); err != nil {
		return nil, fmt.Errorf("create match: %w", err)
	}

	return &SwipeResult{Matched: true, Match: match}, nil
}
