package repository

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, age, language, discord, telegram, region)
		VALUES (:username, :email, :password_hash, :age, :language, :discord, :telegram, :region)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	}
	return fmt.Errorf("create user: no rows returned")
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}
