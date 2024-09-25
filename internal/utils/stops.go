package utils

import (
	"math"
	"sort"

	"github.com/eryalito/vigo-bus-core/pkg/api"
)

func SortStopsByDistance(lat, lon float64, stops []api.Stop) {
	sort.Slice(stops, func(i, j int) bool {
		distI := haversine(lat, lon, float64(stops[i].Location.Lat), float64(stops[i].Location.Lon))
		distJ := haversine(lat, lon, float64(stops[j].Location.Lat), float64(stops[j].Location.Lon))
		return distI < distJ
	})
}

// Haversine formula to calculate the distance between two points on the Earth's surface
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
