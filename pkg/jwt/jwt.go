package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Token struct {
	Email            string
	SignedKey        string
	ExpiresAtMinutes uint
}

func New(email string, signedKey string, expiresAtMinutes uint) *Token {
	return &Token{Email: email, SignedKey: signedKey, ExpiresAtMinutes: expiresAtMinutes}
}

// Generate - метод получения токена. На выходе токен
func (t *Token) Generate() (string, error) {
	duration := time.Duration(t.ExpiresAtMinutes)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tokenString, err := token.SignedString([]byte(t.SignedKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// IsValid - верификация токена
func (t Token) IsValid(token string) bool {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return t.SignedKey, nil
	})
	if err != nil || parsedToken == nil {
		return false
	}

	return parsedToken.Valid
}
