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

func New(email string, expiresAtMinutes uint) Token {
	duration := time.Duration(expiresAtMinutes)
	return Token{Email: email, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * duration)),
	}}
}

// Generate - метод получения токена.
func (t Token) Generate(signedKey string) (string, error) {
	//duration := time.Duration(t.ExpiresAtMinutes)
	//claims := &jwt.RegisteredClaims{
	//	ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration * time.Minute)),
	//	IssuedAt:  jwt.NewNumericDate(time.Now()),
	//	NotBefore: jwt.NewNumericDate(time.Now()),
	//}
	//token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	//tokenString, err := token.SignedString([]byte(t.SignedKey))
	//if err != nil {
	//	return "", err
	//}
	//claims := jwt.MapClaims{}
	//claims["authorized"] = true
	//claims["email"] = t.Email
	//claims["exp"] = time.Now().Add(time.Minute * duration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	tokenString, err := token.SignedString([]byte(signedKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// IsValid - верификация токена
func IsValid(token string, signedKey string) bool {
	//if strings.Contains(token, "Bearer") {
	//	token = strings.Replace(token, "Bearer", "", 1)
	//	token = strings.TrimSpace(token)
	//}
	//parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("There was an error in parsing")
	//	}
	//	return signedKey, nil
	//},
	//)
	//if err != nil || parsedToken == nil {
	//	return false
	//}

	//return parsedToken.Valid
	return validateToken(token, signedKey)
}

func validateToken(token string, signedKey string) bool {
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
