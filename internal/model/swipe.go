package model

import "time"

const (
	SwipeActionLike    = "like"
	SwipeActionDislike = "dislike"
)

type Swipe struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	TargetCardID int64     `db:"target_card_id" json:"target_card_id"`
	Action       string    `db:"action" json:"action"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
