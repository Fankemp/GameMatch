package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	sigingKey string
	ttl       time.Duration
}

func NewManager(singingKey string, ttl time.Duration) *Manager {
	return &Manager{
		sigingKey: singingKey,
		ttl:       ttl,
	}
}

func (m Manager) NewJWT(id int64) (string, error) {
	claims := &tokenClaims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.sigingKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func (m Manager) Parse(tokenStr string) (int64, error) {
	claims := tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(m.sigingKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}
