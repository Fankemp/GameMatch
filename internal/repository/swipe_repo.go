package repository

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/jmoiron/sqlx"
)

type SwipeRepository interface {
	Create(ctx context.Context, swipe *model.Swipe) error
	Exists(ctx context.Context, userID, targetCardID int64) (bool, error)
	CheckMutualLike(ctx context.Context, targetUserID, currentUserID int64, gameID string) (bool, error)
}

type swipeRepo struct {
	db *sqlx.DB
}

func NewSwipeRepository(db *sqlx.DB) SwipeRepository {
	return &swipeRepo{db: db}
}

func (r *swipeRepo) Create(ctx context.Context, swipe *model.Swipe) error {
	query := `
		INSERT INTO swipes (user_id, target_card_id, action)
		VALUES (:user_id, :target_card_id, :action)
		RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, swipe)
	if err != nil {
		return fmt.Errorf("create swipe: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&swipe.ID, &swipe.CreatedAt)
	}
	return fmt.Errorf("create swipe: no rows returned")
}

func (r *swipeRepo) Exists(ctx context.Context, userID, targetCardID int64) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM swipes WHERE user_id = $1 AND target_card_id = $2)`,
		userID, targetCardID)
	return exists, err
}

// CheckMutualLike checks if targetUserID has liked any card of currentUserID for the same game
func (r *swipeRepo) CheckMutualLike(ctx context.Context, targetUserID, currentUserID int64, gameID string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS(
			SELECT 1 FROM swipes s
			JOIN game_cards gc ON gc.id = s.target_card_id
			WHERE s.user_id = $1
			  AND gc.user_id = $2
			  AND gc.game_id = $3
			  AND s.action = 'like'
		)`, targetUserID, currentUserID, gameID)
	return exists, err
}
