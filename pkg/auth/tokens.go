package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	singingKey string
	ttl        time.Duration
}

func (t TokenManager) NewJWT(id int64) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}
