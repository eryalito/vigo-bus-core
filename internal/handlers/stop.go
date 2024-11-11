package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/eryalito/vigo-bus-core/internal/config"
	"github.com/eryalito/vigo-bus-core/internal/sqlite"
	"github.com/eryalito/vigo-bus-core/internal/utils"
	"github.com/eryalito/vigo-bus-core/internal/vitrasa"
	"github.com/eryalito/vigo-bus-core/pkg/api"

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

// GetStopSchedule godoc
// @Summary Get the schedule for a stop
// @Description Provide the schedule for a stop
// @Tags Bus
// @Produce  json
// @Param stop_number path int true "Stop Number"
// @Success 200 {object} api.StopSchedule
// @Router /api/stops/{stop_number}/schedule [get]
func GetStopSchedule(c *gin.Context) {
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

	if stop == (api.Stop{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		return
	}

	vitrasa_client := vitrasa.NewVitrasaClient()
	schedule, err := vitrasa_client.GetSchedules(stop.StopNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &api.StopSchedule{
		Stop:      stop,
		Schedules: schedule,
	})
}

// GetNearbyStopsImage godoc
// @Summary Get the nearby stops as a PNG image and JSON array
// @Description Provide the nearby stops for a location and return a PNG image and JSON array
// @Tags Bus
// @Produce  json
// @Param lat query float64 true "Latitude"
// @Param lon query float64 true "Longitude"
// @Param radius query float64 true "Radius in meters"
// @Param limit query int false "Limit of stops to return, default 9"
// @Success 200 {object} api.NearbyStops
// @Router /api/stops/find/location/image [get]
func GetNearbyStopsImage(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(c.Query("radius"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "9"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
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

	// Sort the stops by distance
	utils.SortStopsByDistance(lat, lon, stops)

	// Truncate the stops to the limit
	if len(stops) > limit {
		stops = stops[:limit]
	}

	// Create the image
	img, err := utils.GenerateImageWithMarkers(config.GoogleMapsAPIKey, struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}{
		Lat: lat,
		Lon: lon,
	}, stops)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image"})
		return
	}

	// Encode the image to base64
	encodedImage, err := utils.PngToBase64(img)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image"})
		return
	}

	// Create the NearbyStops object
	nearbyStops := api.NearbyStops{
		Stops:   stops,
		Image64: encodedImage,
		Radius:  radius,
		Origin: struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		}{
			Lat: lat,
			Lon: lon,
		},
	}

	// Return the NearbyStops object
	c.JSON(http.StatusOK, nearbyStops)
}
