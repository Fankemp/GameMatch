package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func (b PasswordHasher) Hash(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), b.cost)
	return string(hash), err
}

func (b PasswordHasher) Compare()
