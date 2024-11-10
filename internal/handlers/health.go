package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListLines godoc
// @Summary Health endpoint
// @Description Health endpoint
// @Tags Health
// @Produce text/plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
