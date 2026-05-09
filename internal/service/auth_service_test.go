package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/service"
	"github.com/Fankemp/GameMatch/internal/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestAuthService_SignUp(t *testing.T) {
	type mockBehavior func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignUpInput)

	tests := []struct {
		name         string
		input        service.SignUpInput
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "Success",
			input: service.SignUpInput{
				Username: "Sarvan",
				Email:    "musaevsarvan960@gmail.com",
				Password: "password123",
			},
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignUpInput) {
				r.EXPECT().GetByEmail(gomock.Any(), input.Email).Return(nil, model.ErrUserNotFound)

				h.EXPECT().Hash(input.Password).Return("hashed_pw", nil)

				r.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, u *model.User) error {
					u.ID = 1
					return nil
				})

				tm.EXPECT().NewJWT(int64(1)).Return("token123", nil)
			},
			wantErr: nil,
		},
		{
			name: "User already exists",
			input: service.SignUpInput{
				Email: "exist@test.com",
			},
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignUpInput) {
				r.EXPECT().GetByEmail(gomock.Any(), input.Email).Return(&model.User{ID: 1}, nil)
			},
			wantErr: service.ErrUserAlreadyExists,
		},
		{
			name:  "Database Error",
			input: service.SignUpInput{Email: "db_error@test.com"},
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignUpInput) {
				r.EXPECT().GetByEmail(gomock.Any(), input.Email).Return(nil, service.ErrInternalDB)
			},
			wantErr: service.ErrInternalDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mocks.NewMockUserRepository(ctrl)
			h := mocks.NewMockPasswordHasher(ctrl)
			tm := mocks.NewMockTokenManager(ctrl)

			svc := service.NewAuthService(r, h, tm)
			tt.mockBehavior(r, h, tm, tt.input)
			resp, err := svc.SignUp(context.Background(), tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}

				if resp.Token == "" {
					t.Error("expected  non empty token")
				}

				if resp.User.PasswordHash != "" {
					t.Error("security breach: password hash returned")
				}
			}

		})
	}
}

func TestAuthService_SignIn(t *testing.T) {
	type mockBehavior func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignInInput)

	tests := []struct {
		name         string
		input        service.SignInInput
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "Success",
			input: service.SignInInput{
				Email:    "success@test.com",
				Password: "password123",
			},
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignInInput) {
				userInDB := &model.User{
					ID:           1,
					Email:        input.Email,
					PasswordHash: "any_hash",
				}
				r.EXPECT().GetByEmail(gomock.Any(), input.Email).Return(userInDB, nil)
				h.EXPECT().Compare(userInDB.PasswordHash, input.Password).Return(nil)
				tm.EXPECT().NewJWT(int64(1)).Return("new_token", nil)
			},
			wantErr: nil,
		},
		{
			name: "User not exist",
			input: service.SignInInput{
				Email:    "existsuser@test.com",
				Password: "password123",
			},
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignInInput) {
				r.EXPECT().GetByEmail(gomock.Any(), input.Email).Return(nil, model.ErrUserNotFound)
			},
			wantErr: service.ErrInvalidCredentials,
		},
		{
			name: "Invalid password",
			input: service.SignInInput{
				Email:    "test@test.com",
				Password: "pasword123",
			},
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input service.SignInInput) {
				userInDb := &model.User{
					ID:           1,
					PasswordHash: "hash123",
				}

				r.EXPECT().GetByEmail(gomock.Any(), input.Email).Return(userInDb, nil)
				h.EXPECT().Compare(userInDb.PasswordHash, input.Password).Return(errors.New("Invalid password"))
			},
			wantErr: service.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := mocks.NewMockUserRepository(ctrl)
			h := mocks.NewMockPasswordHasher(ctrl)
			tm := mocks.NewMockTokenManager(ctrl)

			svc := service.NewAuthService(r, h, tm)
			tt.mockBehavior(r, h, tm, tt.input)
			resp, err := svc.SignIn(context.Background(), tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}

				if resp.Token == "" {
					t.Error("expected  non empty token")
				}

				if resp.User.PasswordHash != "" {
					t.Error("security breach: password hash returned")
				}
			}

		})
	}
}
func TestAuthService_GetMe(t *testing.T) {
	type mockBehavior func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input int64)

	tests := []struct {
		name         string
		input        int64
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name:  "Success",
			input: 1,
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input int64) {
				r.EXPECT().GetByID(gomock.Any(), input).Return(&model.User{ID: 1}, nil)
			},
			wantErr: nil,
		},
		{
			name:  "Not Found",
			input: 2,
			mockBehavior: func(r *mocks.MockUserRepository, h *mocks.MockPasswordHasher, tm *mocks.MockTokenManager, input int64) {
				r.EXPECT().GetByID(gomock.Any(), input).Return(nil, model.ErrUserNotFound)
			},
			wantErr: model.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := mocks.NewMockUserRepository(ctrl)
			h := mocks.NewMockPasswordHasher(ctrl)
			tm := mocks.NewMockTokenManager(ctrl)

			svc := service.NewAuthService(r, h, tm)
			tt.mockBehavior(r, h, tm, tt.input)
			resp, err := svc.GetMe(context.Background(), tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error: %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil {
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.PasswordHash != "" {
					t.Error("security breach: password hash returned in GetMe")
				}
			}
		})
	}
}
