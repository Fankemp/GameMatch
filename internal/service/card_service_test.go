package service

import (
	"context"
	"testing"

	"github.com/Fankemp/GameMatch/internal/repository/mocks"
)

func TestCardCreate_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepo()
	svc := NewCardService(cardRepo)

	input := CreateCardInput{
		GameID:      "valorant",
		Rank:        "Diamond",
		Role:        "Duelist",
		Description: "Jett main, looking for duo",
	}

	card, err := svc.Create(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if card.ID == 0 {
		t.Fatal("expected card ID to be set")
	}
	if card.UserID != 1 {
		t.Fatalf("expected user_id 1, got %d", card.UserID)
	}
	if card.GameID != "valorant" {
		t.Fatalf("expected game_id 'valorant', got '%s'", card.GameID)
	}
	if !card.IsActive {
		t.Fatal("expected card to be active by default")
	}
}

func TestCardGetMyCards_Empty(t *testing.T) {
	cardRepo := mocks.NewCardRepo()
	svc := NewCardService(cardRepo)

	cards, err := svc.GetMyCards(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cards) != 0 {
		t.Fatalf("expected 0 cards, got %d", len(cards))
	}
}

func TestCardGetMyCards_ReturnsOwned(t *testing.T) {
	cardRepo := mocks.NewCardRepo()
	svc := NewCardService(cardRepo)

	svc.Create(context.Background(), 1, CreateCardInput{GameID: "valorant", Rank: "Gold", Role: "Duelist", Description: "test"})
	svc.Create(context.Background(), 1, CreateCardInput{GameID: "cs2", Rank: "Silver", Role: "Sentinel", Description: "test2"})
	svc.Create(context.Background(), 2, CreateCardInput{GameID: "valorant", Rank: "Iron", Role: "Controller", Description: "other"})

	cards, err := svc.GetMyCards(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cards) != 2 {
		t.Fatalf("expected 2 cards for user 1, got %d", len(cards))
	}
}

func TestCardUpdate_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepo()
	svc := NewCardService(cardRepo)

	svc.Create(context.Background(), 1, CreateCardInput{GameID: "valorant", Rank: "Gold", Role: "Duelist", Description: "old"})

	isActive := false
	updated, err := svc.Update(context.Background(), 1, 1, UpdateCardInput{
		Rank:     "Diamond",
		IsActive: &isActive,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Rank != "Diamond" {
		t.Fatalf("expected rank 'Diamond', got '%s'", updated.Rank)
	}
	if updated.IsActive {
		t.Fatal("expected card to be inactive")
	}
}

func TestCardUpdate_WrongOwner(t *testing.T) {
	cardRepo := mocks.NewCardRepo()
	svc := NewCardService(cardRepo)

	svc.Create(context.Background(), 1, CreateCardInput{GameID: "valorant", Rank: "Gold", Role: "Duelist", Description: "mine"})

	_, err := svc.Update(context.Background(), 2, 1, UpdateCardInput{Rank: "Iron"})
	if err == nil {
		t.Fatal("expected error when updating someone else's card")
	}
}

func TestCardDelete_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepo()
	svc := NewCardService(cardRepo)

	svc.Create(context.Background(), 1, CreateCardInput{GameID: "valorant", Rank: "Gold", Role: "Duelist", Description: "to delete"})

	err := svc.Delete(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cards, _ := svc.GetMyCards(context.Background(), 1)
	if len(cards) != 0 {
		t.Fatalf("expected 0 cards after delete, got %d", len(cards))
	}
}

func TestCardDelete_WrongOwner(t *testing.T) {
	cardRepo := mocks.NewCardRepo()
	svc := NewCardService(cardRepo)

	svc.Create(context.Background(), 1, CreateCardInput{GameID: "valorant", Rank: "Gold", Role: "Duelist", Description: "mine"})

	err := svc.Delete(context.Background(), 2, 1)
	if err == nil {
		t.Fatal("expected error when deleting someone else's card")
	}
}
