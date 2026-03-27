package model

import "time"

type Match struct {
	ID        int64     `db:"id" json:"id"`
	UserAID   int64     `db:"user_a_id" json:"user_a_id"`
	UserBID   int64     `db:"user_b_id" json:"user_b_id"`
	GameID    string    `db:"game_id" json:"game_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
