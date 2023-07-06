package middleware

import (
	"github.com/cost_control/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func AuthMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !jwt.IsValid(token, os.Getenv("SIGNED_KEY")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()

}
