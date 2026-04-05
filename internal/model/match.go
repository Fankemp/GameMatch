package model

import "time"

type Match struct {
	ID        int64     `db:"id" json:"id"`
	UserAID   int64     `db:"user_a_id" json:"user_a_id"`
	UserBID   int64     `db:"user_b_id" json:"user_b_id"`
	GameID    string    `db:"game_id" json:"game_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type MatchUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
	Language string `json:"language"`
	Region   string `json:"region"`
	Discord  string `json:"discord,omitempty"`
	Telegram string `json:"telegram,omitempty"`
}

type MatchWithUser struct {
	Match Match     `json:"match"`
	User  MatchUser `json:"user"`
	Rank  string    `json:"rank"`
	Role  string    `json:"role"`
}
