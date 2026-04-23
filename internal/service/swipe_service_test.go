package service

import (
	"context"
	"testing"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository/mocks"
)

func setupSwipeTest() (*mocks.SwipeRepo, *mocks.CardRepo, *mocks.MatchRepo, SwipeService) {
	swipeRepo := mocks.NewSwipeRepo()
	cardRepo := mocks.NewCardRepo()
	matchRepo := mocks.NewMatchRepo()

	// Create a card owned by user 2
	cardRepo.Cards[1] = &model.GameCard{
		ID: 1, UserID: 2, GameID: "valorant", Rank: "Diamond", Role: "Duelist", IsActive: true,
	}
	cardRepo.NextID = 2

	svc := NewSwipeService(swipeRepo, cardRepo, matchRepo)
	return swipeRepo, cardRepo, matchRepo, svc
}

func TestSwipe_Like_NoMutual(t *testing.T) {
	_, _, _, svc := setupSwipeTest()

	result, err := svc.Swipe(context.Background(), 1, SwipeInput{
		TargetCardID: 1,
		Action:       "like",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Matched {
		t.Fatal("expected no match when no mutual like")
	}
	if result.Match != nil {
		t.Fatal("expected match to be nil")
	}
}

func TestSwipe_Like_MutualMatch(t *testing.T) {
	swipeRepo, _, matchRepo, svc := setupSwipeTest()

	// Simulate that user 2 already liked user 1's card
	swipeRepo.MutualLike = true

	result, err := svc.Swipe(context.Background(), 1, SwipeInput{
		TargetCardID: 1,
		Action:       "like",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Matched {
		t.Fatal("expected match when mutual like exists")
	}
	if result.Match == nil {
		t.Fatal("expected match object to be set")
	}
	if result.Match.UserAID != 1 || result.Match.UserBID != 2 {
		t.Fatalf("expected match between user 1 and 2, got %d and %d", result.Match.UserAID, result.Match.UserBID)
	}
	if result.Match.GameID != "valorant" {
		t.Fatalf("expected game_id 'valorant', got '%s'", result.Match.GameID)
	}
	if len(matchRepo.Matches) != 1 {
		t.Fatalf("expected 1 match in repo, got %d", len(matchRepo.Matches))
	}
}

func TestSwipe_Dislike_NoMatch(t *testing.T) {
	swipeRepo, _, _, svc := setupSwipeTest()
	swipeRepo.MutualLike = true // even if mutual exists

	result, err := svc.Swipe(context.Background(), 1, SwipeInput{
		TargetCardID: 1,
		Action:       "dislike",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Matched {
		t.Fatal("dislike should never create a match")
	}
}

func TestSwipe_DuplicateSwipe(t *testing.T) {
	_, _, _, svc := setupSwipeTest()

	// First swipe
	_, err := svc.Swipe(context.Background(), 1, SwipeInput{TargetCardID: 1, Action: "like"})
	if err != nil {
		t.Fatalf("first swipe failed: %v", err)
	}

	// Second swipe on same card — should silently return no match
	result, err := svc.Swipe(context.Background(), 1, SwipeInput{TargetCardID: 1, Action: "like"})
	if err != nil {
		t.Fatalf("expected no error on duplicate, got %v", err)
	}
	if result.Matched {
		t.Fatal("duplicate swipe should not create match")
	}
}

func TestSwipe_OwnCard(t *testing.T) {
	swipeRepo := mocks.NewSwipeRepo()
	cardRepo := mocks.NewCardRepo()
	matchRepo := mocks.NewMatchRepo()

	// Card owned by user 1
	cardRepo.Cards[1] = &model.GameCard{
		ID: 1, UserID: 1, GameID: "valorant", Rank: "Gold", Role: "Duelist", IsActive: true,
	}

	svc := NewSwipeService(swipeRepo, cardRepo, matchRepo)

	_, err := svc.Swipe(context.Background(), 1, SwipeInput{TargetCardID: 1, Action: "like"})
	if err == nil {
		t.Fatal("expected error when swiping on own card")
	}
}

func TestSwipe_MatchAlreadyExists(t *testing.T) {
	swipeRepo, _, matchRepo, svc := setupSwipeTest()
	swipeRepo.MutualLike = true
	matchRepo.MatchExists = true // match already exists

	result, err := svc.Swipe(context.Background(), 1, SwipeInput{TargetCardID: 1, Action: "like"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Matched {
		t.Fatal("should not create duplicate match")
	}
	if len(matchRepo.Matches) != 0 {
		t.Fatalf("expected 0 matches in repo, got %d", len(matchRepo.Matches))
	}
}
