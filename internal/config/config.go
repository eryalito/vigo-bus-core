package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	Port             string
	Token            string
	StopsDBPath      string
	IdentityDBPath   string
	GoogleMapsAPIKey string
	RateLimiter      struct {
		Limit int
		Burst int
	}
)

func Init() {
	// Define command-line flags
	flag.StringVar(&Port, "port", getEnv("PORT", "8080"), "Port to run the server on")
	flag.StringVar(&Token, "token", getEnv("TOKEN", "your-secret-token"), "Authentication token")
	flag.StringVar(&StopsDBPath, "stops-db-path", getEnv("STOPS_DB_PATH", "stops.db"), "Path to the stops database")
	flag.StringVar(&IdentityDBPath, "identity-db-path", getEnv("IDENTITY_DB_PATH", "identity.db"), "Path to the identity database")
	flag.StringVar(&GoogleMapsAPIKey, "google-maps-api-key", getEnv("GOOGLE_MAPS_API_KEY", ""), "Google maps api key for generating images")
	limit, err := strconv.Atoi(getEnv("RATE_LIMITER_LIMIT", "1"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse RATE_LIMITER_LIMIT: %v", err))
	}
	flag.IntVar(&RateLimiter.Limit, "rate-limiter-limit", limit, "Rate limiter limit")
	burst, err := strconv.Atoi(getEnv("RATE_LIMITER_BURST", "5"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse RATE_LIMITER_BURST: %v", err))
	}
	flag.IntVar(&RateLimiter.Burst, "rate-limiter-burst", burst, "Rate limiter burst")

	// Parse command-line flags
	flag.Parse()
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
