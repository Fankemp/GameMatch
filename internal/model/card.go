package model

import "time"

type GameCard struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	GameID      string    `db:"game_id" json:"game_id"`
	Rank        string    `db:"rank" json:"rank"`
	Role        string    `db:"role" json:"role"`
	Description string    `db:"description" json:"description"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
