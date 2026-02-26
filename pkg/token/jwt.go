package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenManager interface {
	CreateAccessToken(userID string) (string, error)
	CreateRefreshToken(userID string) (string, error)
	VerifyToken(tokenString string) (*Claims, error)
}

type JWTManager struct {
	secretKey             string
	accessExpirationHours int
	refreshExpirationDays int
}

type Claims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, accessExpirationHours int, refreshExpirationDays int) *JWTManager {
	return &JWTManager{
		secretKey:             secretKey,
		accessExpirationHours: accessExpirationHours,
		refreshExpirationDays: refreshExpirationDays,
	}
}

func (m *JWTManager) CreateAccessToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(m.accessExpirationHours) * time.Hour)
	return m.createToken(userID, "access", expirationTime)
}

func (m *JWTManager) CreateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(m.refreshExpirationDays) * 24 * time.Hour)
	return m.createToken(userID, "refresh", expirationTime)
}

func (m *JWTManager) createToken(userID string, tokenType string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
