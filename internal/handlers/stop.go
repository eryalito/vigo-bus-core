package handlers

import (
	"net/http"
	"strconv"

	"github.com/eryalito/vigo-bus-core/internal/sqlite"

	"github.com/gin-gonic/gin"
)

// ListStops godoc
// @Summary List all of the stops
// @Description Provide a list of all the stops
// @Tags Bus
// @Produce  json
// @Success 200 {array} api.Stop
// @Router /api/stops [get]
func ListStops(c *gin.Context) {
	sdb_conn, err := sqlite.NewBusConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stops, err := sdb_conn.GetStops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stops)
}

// GetStop godoc
// @Summary Get a stop by its number
// @Description Provide a stop by its number
// @Tags Bus
// @Produce  json
// @Param stop_number path int true "Stop Number"
// @Success 200 {object} api.Stop
// @Router /api/stops/{stop_number} [get]
func GetStop(c *gin.Context) {
	stopNumber := c.Param("stop_number")
	stopNumberInt, err := strconv.Atoi(stopNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stop number"})
		return
	}

	sdb_conn, err := sqlite.NewBusConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stop, err := sdb_conn.GetStopByNumber(stopNumberInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stop)
}

// FindStop godoc
// @Summary Find a stop by text in its name
// @Description Provide a list of stops that match the text in their name
// @Tags Bus
// @Produce  json
// @Param text query string true "Text to search for in stop name"
// @Success 200 {array} api.Stop
// @Router /api/stops/find [get]
func FindStops(c *gin.Context) {
	text := c.Query("text")
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing text query parameter"})
		return
	}

	sdb_conn, err := sqlite.NewBusConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stops, err := sdb_conn.FindStopsByText(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stops)
}

// FindStopsByLocation godoc
// @Summary Find a stop by its location
// @Description Provide a list of stops in a given radius around a location
// @Tags Bus
// @Produce  json
// @Param lat query float64 true "Latitude"
// @Param lon query float64 true "Longitude"
// @Param radius query float64 true "Radius in meters"
// @Success 200 {array} api.Stop
// @Router /api/stops/find/location [get]
func FindStopsByLocation(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.Query("radius")
	if latStr == "" || lonStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing lat, lon, or radius query parameters"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lat query parameter"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lon query parameter"})
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius query parameter"})
		return
	}

	sdb_conn, err := sqlite.NewBusConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stops, err := sdb_conn.FindStopsByLocation(lat, lon, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stops)
}
