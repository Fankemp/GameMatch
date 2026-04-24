package model

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Age          int       `db:"age" json:"age"`
	Language     string    `db:"language" json:"language"`
	Discord      string    `db:"discord" json:"discord,omitempty"`
	Telegram     string    `db:"telegram" json:"telegram,omitempty"`
	Region       string    `db:"region" json:"region"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
