package middleware

import (
	"net/http"
	"strings"

	"github.com/eryalito/vigo-bus-core/internal/config"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if config.Token != token {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}
