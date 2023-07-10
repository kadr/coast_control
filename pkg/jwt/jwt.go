package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Token struct {
	Email string
	jwt.RegisteredClaims
}

func New() *Token {
	return &Token{}
}

// Generate - метод получения токена.
func (t Token) Generate(email string, expiresAtMinutes uint, signedKey string) (string, error) {
	duration := time.Duration(expiresAtMinutes)
	t.Email = email
	t.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * duration))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	tokenString, err := token.SignedString([]byte(signedKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// IsValid - верификация токена
func (t Token) IsValid(token string, signedKey string) bool {
	parsedToken, err := parseToken(token, signedKey)
	if err != nil {
		return false
	}
	if _, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return true
	}

	return false
}

func parseToken(token string, signedKey string) (*jwt.Token, error) {
	if strings.Contains(token, "Bearer") {
		token = strings.Replace(token, "Bearer", "", 1)
		token = strings.TrimSpace(token)
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(signedKey), nil
	})
	return parsedToken, err
}
