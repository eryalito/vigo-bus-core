package api

type NearbyStops struct {
	Origin struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"origin"`
	Radius  float64 `json:"radius"`
	Stops   []Stop  `json:"stops"`
	Image64 string  `json:"image"`
}
