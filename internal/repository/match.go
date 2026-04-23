package repository

import (
	"context"
	"fmt"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/jmoiron/sqlx"
)

type matchRepo struct {
	db *sqlx.DB
}

func NewMatchRepository(db *sqlx.DB) MatchRepository {
	return &matchRepo{db: db}
}

func (r *matchRepo) Create(ctx context.Context, match *model.Match) error {
	query := `
		INSERT INTO matches (user_a_id, user_b_id, game_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		match.UserAID, match.UserBID, match.GameID,
	).Scan(&match.ID, &match.CreatedAt)
}

func (r *matchRepo) GetByUserID(ctx context.Context, userID int64) ([]*model.MatchWithUser, error) {
	query := `
		SELECT m.id AS match_id, m.user_a_id, m.user_b_id, m.game_id, m.created_at AS match_created_at,
		       u.id AS user_id, u.username, u.age, u.language, u.region, u.discord, u.telegram,
		       gc.rank, gc.role
		FROM matches m
		JOIN users u ON u.id = CASE WHEN m.user_a_id = $1 THEN m.user_b_id ELSE m.user_a_id END
		LEFT JOIN game_cards gc ON gc.user_id = u.id AND gc.game_id = m.game_id AND gc.is_active = true
		WHERE m.user_a_id = $1 OR m.user_b_id = $1
		ORDER BY m.created_at DESC`

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get matches: %w", err)
	}
	defer rows.Close()

	var matches []*model.MatchWithUser
	for rows.Next() {
		var mwu model.MatchWithUser
		err := rows.Scan(
			&mwu.Match.ID, &mwu.Match.UserAID, &mwu.Match.UserBID,
			&mwu.Match.GameID, &mwu.Match.CreatedAt,
			&mwu.User.ID, &mwu.User.Username, &mwu.User.Age,
			&mwu.User.Language, &mwu.User.Region, &mwu.User.Discord, &mwu.User.Telegram,
			&mwu.Rank, &mwu.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("scan match: %w", err)
		}
		matches = append(matches, &mwu)
	}
	return matches, nil
}

func (r *matchRepo) Exists(ctx context.Context, userAID, userBID int64, gameID string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `
		SELECT EXISTS(
			SELECT 1 FROM matches
			WHERE LEAST(user_a_id, user_b_id) = LEAST($1, $2)
			  AND GREATEST(user_a_id, user_b_id) = GREATEST($1, $2)
			  AND game_id = $3
		)`, userAID, userBID, gameID)
	return exists, err
}
