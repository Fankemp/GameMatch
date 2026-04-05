package repository

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/jmoiron/sqlx"
)

type CardRepository interface {
	Create(ctx context.Context, card *model.GameCard) error
	GetByID(ctx context.Context, id int64) (*model.GameCard, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.GameCard, error)
	Update(ctx context.Context, card *model.GameCard) error
	Delete(ctx context.Context, id, userID int64) error
	GetFeed(ctx context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error)
}

type cardRepo struct {
	db *sqlx.DB
}

func NewCardRepository(db *sqlx.DB) CardRepository {
	return &cardRepo{db: db}
}

func (r *cardRepo) Create(ctx context.Context, card *model.GameCard) error {
	query := `
		INSERT INTO game_cards (user_id, game_id, rank, role, description)
		VALUES (:user_id, :game_id, :rank, :role, :description)
		RETURNING id, is_active, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, card)
	if err != nil {
		return fmt.Errorf("create card: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&card.ID, &card.IsActive, &card.CreatedAt, &card.UpdatedAt)
	}
	return fmt.Errorf("create card: no rows returned")
}

func (r *cardRepo) GetByID(ctx context.Context, id int64) (*model.GameCard, error) {
	var card model.GameCard
	err := r.db.GetContext(ctx, &card, `SELECT * FROM game_cards WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	return &card, nil
}

func (r *cardRepo) GetByUserID(ctx context.Context, userID int64) ([]*model.GameCard, error) {
	var cards []*model.GameCard
	err := r.db.SelectContext(ctx, &cards,
		`SELECT * FROM game_cards WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("get cards by user: %w", err)
	}
	return cards, nil
}

func (r *cardRepo) Update(ctx context.Context, card *model.GameCard) error {
	query := `
		UPDATE game_cards
		SET rank = $1, role = $2, description = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		card.Rank, card.Role, card.Description, card.IsActive, card.ID, card.UserID,
	).Scan(&card.UpdatedAt)
}

func (r *cardRepo) Delete(ctx context.Context, id, userID int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM game_cards WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("card not found")
	}
	return nil
}

func (r *cardRepo) GetFeed(ctx context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error) {
	query := `
		SELECT gc.id, gc.user_id, gc.game_id, gc.rank, gc.role, gc.description,
		       gc.is_active, gc.created_at, gc.updated_at,
		       u.username, u.age, u.language, u.region, u.discord, u.telegram
		FROM game_cards gc
		JOIN users u ON u.id = gc.user_id
		WHERE gc.game_id = $1
		  AND gc.user_id != $2
		  AND gc.is_active = true
		  AND NOT EXISTS (
		      SELECT 1 FROM swipes s WHERE s.user_id = $2 AND s.target_card_id = gc.id
		  )
		ORDER BY gc.created_at DESC
		LIMIT $3 OFFSET $4`

	var cards []*model.FeedCard
	err := r.db.SelectContext(ctx, &cards, query, gameID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get feed: %w", err)
	}
	return cards, nil
}
