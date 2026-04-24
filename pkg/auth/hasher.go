package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{
		cost: cost,
	}
}

func (b BcryptHasher) Hash(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), b.cost)
	return string(hash), err
}

func (b BcryptHasher) Compare(userPassword string, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword))
	return err
}
