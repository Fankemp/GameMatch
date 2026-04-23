package service

import (
	"context"
	"testing"

	"github.com/Fankemp/GameMatch/internal/repository/mocks"
)

const testJWTSecret = "test-secret-key"

func TestRegister_Success(t *testing.T) {
	userRepo := mocks.NewUserRepo()
	svc := NewAuthService(userRepo, testJWTSecret)

	input := RegisterInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Age:      20,
		Language: "RU",
		Region:   "CIS",
	}

	resp, err := svc.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token to be non-empty")
	}
	if resp.User.Username != "testuser" {
		t.Fatalf("expected username 'testuser', got '%s'", resp.User.Username)
	}
	if resp.User.Email != "test@example.com" {
		t.Fatalf("expected email 'test@example.com', got '%s'", resp.User.Email)
	}
	if resp.User.ID == 0 {
		t.Fatal("expected user ID to be set")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	userRepo := mocks.NewUserRepo()
	svc := NewAuthService(userRepo, testJWTSecret)

	input := RegisterInput{
		Username: "user1",
		Email:    "dup@example.com",
		Password: "password123",
		Age:      20,
		Language: "RU",
		Region:   "CIS",
	}

	_, err := svc.Register(context.Background(), input)
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	input.Username = "user2"
	_, err = svc.Register(context.Background(), input)
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
	if err != ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	userRepo := mocks.NewUserRepo()
	svc := NewAuthService(userRepo, testJWTSecret)

	// Register first
	regInput := RegisterInput{
		Username: "logintest",
		Email:    "login@example.com",
		Password: "mypassword",
		Age:      25,
		Language: "EN",
		Region:   "EU Central",
	}
	_, err := svc.Register(context.Background(), regInput)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Login
	loginInput := LoginInput{
		Email:    "login@example.com",
		Password: "mypassword",
	}
	resp, err := svc.Login(context.Background(), loginInput)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token to be non-empty")
	}
	if resp.User.Username != "logintest" {
		t.Fatalf("expected username 'logintest', got '%s'", resp.User.Username)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := mocks.NewUserRepo()
	svc := NewAuthService(userRepo, testJWTSecret)

	regInput := RegisterInput{
		Username: "wrongpass",
		Email:    "wrong@example.com",
		Password: "correctpassword",
		Age:      20,
		Language: "RU",
		Region:   "CIS",
	}
	_, _ = svc.Register(context.Background(), regInput)

	loginInput := LoginInput{
		Email:    "wrong@example.com",
		Password: "incorrectpassword",
	}
	_, err := svc.Login(context.Background(), loginInput)
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_NonExistentEmail(t *testing.T) {
	userRepo := mocks.NewUserRepo()
	svc := NewAuthService(userRepo, testJWTSecret)

	loginInput := LoginInput{
		Email:    "noone@example.com",
		Password: "password",
	}
	_, err := svc.Login(context.Background(), loginInput)
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestGetMe_Success(t *testing.T) {
	userRepo := mocks.NewUserRepo()
	svc := NewAuthService(userRepo, testJWTSecret)

	regInput := RegisterInput{
		Username: "meuser",
		Email:    "me@example.com",
		Password: "password123",
		Age:      22,
		Language: "KZ",
		Region:   "CIS",
	}
	resp, _ := svc.Register(context.Background(), regInput)

	user, err := svc.GetMe(context.Background(), resp.User.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Username != "meuser" {
		t.Fatalf("expected 'meuser', got '%s'", user.Username)
	}
}
