package repository

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/jmoiron/sqlx"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *model.Profile) error
	GetByUserID(ctx context.Context, userID int64) (*model.Profile, error)
	GetByID(ctx context.Context, id int64) (*model.Profile, error)
	Update(ctx context.Context, profile *model.Profile) error
}

type profileRepo struct {
	db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return &profileRepo{db: db}
}

func (r *profileRepo) Create(ctx context.Context, profile *model.Profile) error {
	query := `
		INSERT INTO profiles (user_id, bio, avatar_url)
		VALUES (:user_id, :bio, :avatar_url)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, profile)
	if err != nil {
		return fmt.Errorf("create profile: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
	}
	return fmt.Errorf("create profile: no rows returned")
}

func (r *profileRepo) GetByUserID(ctx context.Context, userID int64) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.GetContext(ctx, &profile, `SELECT * FROM profiles WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile by user_id: %w", err)
	}
	return &profile, nil
}

func (r *profileRepo) GetByID(ctx context.Context, id int64) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.GetContext(ctx, &profile, `SELECT * FROM profiles WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get profile by id: %w", err)
	}
	return &profile, nil
}

func (r *profileRepo) Update(ctx context.Context, profile *model.Profile) error {
	query := `
		UPDATE profiles
		SET bio = :bio, avatar_url = :avatar_url, updated_at = NOW()
		WHERE user_id = :user_id
		RETURNING updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, profile)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&profile.UpdatedAt)
	}
	return fmt.Errorf("update profile: no rows returned")
}
