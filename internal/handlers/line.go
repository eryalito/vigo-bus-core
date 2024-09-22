package handlers

import (
	"net/http"

	"github.com/eryalito/vigo-bus-core/internal/sqlite"

	"github.com/gin-gonic/gin"
)

// ListLines godoc
// @Summary List all of the lines
// @Description Provide a list of all the lines
// @Tags Bus
// @Produce  json
// @Success 200 {array} api.Line
// @Router /api/lines [get]
func ListLines(c *gin.Context) {
	sdb_conn, err := sqlite.NewBusConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lines, err := sdb_conn.GetLines()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lines)
}
