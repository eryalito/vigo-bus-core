package handlers

import (
	"net/http"
	"strconv"

	"github.com/eryalito/vigo-bus-core/internal/sqlite"
	"github.com/eryalito/vigo-bus-core/pkg/api"

	"github.com/gin-gonic/gin"
)

// GetUser godoc
// @Summary Get a user by its UUID for a specific provider
// @Description Provide a user by its UUID for a specific provider
// @Tags Identity
// @Produce  json
// @Param provider path string true "Provider"
// @Param uuid path string true "UUID"
// @Success 200 {object} api.Identity
// @Router /api/users/{provider}/{uuid} [get]
func GetUser(c *gin.Context) {
	provider := c.Param("provider")
	uuid := c.Param("uuid")

	sdb_conn, err := sqlite.NewIdentityConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := sdb_conn.GetUserByUUID(provider, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags Identity
// @Produce  json
// @Param provider path string true "Provider"
// @Param uuid path string true "UUID"
// @Success 200 {object} api.Identity
// @Router /api/users/{provider}/{uuid} [post]
func CreateUser(c *gin.Context) {
	provider := c.Param("provider")
	uuid := c.Param("uuid")

	sdb_conn, err := sqlite.NewIdentityConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer sdb_conn.Close()

	// Check if a user with the same UUID and provider already exists
	existingUser, err := sdb_conn.GetUserByUUID(provider, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with the same UUID and provider already exists"})
		return
	}

	identity := &api.Identity{
		Provider: api.ProviderType(provider),
		UUID:     uuid,
	}

	err = sdb_conn.InsertIdentity(identity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := sdb_conn.GetUserByUUID(provider, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// AddFavoriteStopToIdentity godoc
// @Summary Add a favorite stop to a user
// @Description Add a favorite stop to a user
// @Tags Identity
// @Produce  json
// @Param provider path string true "Provider"
// @Param uuid path string true "UUID"
// @Param stop_number path int true "Stop Number"
// @Success 200 {object} api.Identity
// @Router /api/users/{provider}/{uuid}/favorite_stops/{stop_number} [post]
func AddFavoriteStopToIdentity(c *gin.Context) {
	provider := c.Param("provider")
	uuid := c.Param("uuid")
	stopNumber := c.Param("stop_number")

	stopNumberInt, err := strconv.Atoi(stopNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stop number"})
		return
	}

	bdb_conn, err := sqlite.NewBusConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer bdb_conn.Close()

	stop, err := bdb_conn.GetStopByNumber(stopNumberInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if stop == (api.Stop{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stop not found"})
		return
	}

	sdb_conn, err := sqlite.NewIdentityConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer sdb_conn.Close()

	user, err := sdb_conn.GetUserByUUID(provider, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the stop is already a favorite
	for _, stop := range user.FavoriteStops {
		if stop.StopNumber == stopNumberInt {
			c.JSON(http.StatusConflict, gin.H{"error": "Stop is already a favorite"})
			return
		}
	}

	user.FavoriteStops = append(user.FavoriteStops, stop)

	err = sdb_conn.UpdateIdentity(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RemoveFavoriteStopFromIdentity godoc
// @Summary Remove a favorite stop from a user
// @Description Remove a favorite stop from a user
// @Tags Identity
// @Produce  json
// @Param provider path string true "Provider"
// @Param uuid path string true "UUID"
// @Param stop_number path int true "Stop Number"
// @Success 200 {object} api.Identity
// @Router /api/users/{provider}/{uuid}/favorite_stops/{stop_number} [delete]
func RemoveFavoriteStopFromIdentity(c *gin.Context) {
	provider := c.Param("provider")
	uuid := c.Param("uuid")
	stopNumber := c.Param("stop_number")

	stopNumberInt, err := strconv.Atoi(stopNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stop number"})
		return
	}

	sdb_conn, err := sqlite.NewIdentityConnector()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer sdb_conn.Close()

	user, err := sdb_conn.GetUserByUUID(provider, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the stop is a favorite
	var stopIndex = -1
	for i, stop := range user.FavoriteStops {
		if stop.StopNumber == stopNumberInt {
			stopIndex = i
			break
		}
	}
	if stopIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stop is not a favorite"})
		return
	}

	user.FavoriteStops = append(user.FavoriteStops[:stopIndex], user.FavoriteStops[stopIndex+1:]...)

	err = sdb_conn.UpdateIdentity(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
