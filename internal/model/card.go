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

type FeedCard struct {
	GameCard
	Username string `db:"username" json:"username"`
	Age      int    `db:"age" json:"age"`
	Language string `db:"language" json:"language"`
	Region   string `db:"region" json:"region"`
	Discord  string `db:"discord" json:"discord,omitempty"`
	Telegram string `db:"telegram" json:"telegram,omitempty"`
}
