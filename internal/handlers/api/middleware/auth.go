package middleware

import (
	"github.com/cost_control/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type AuthMiddleware struct {
	jwtManager *jwt.Token
}

func New() *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwt.New()}
}

func (am AuthMiddleware) AuthMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !am.jwtManager.IsValid(token, os.Getenv("SIGNED_KEY")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()

}
