package api

// Stop is a struct that holds the information of a bus stop
type Stop struct {
	// ID is the unique identifier of the stop
	ID int `json:"id"`

	// StopNumber is the number of the stop provided by the bus company
	StopNumber int `json:"stop_number"`

	// StopID is the number of the stop used internally by the bus company
	StopID int `json:"stop_id"`

	// Name is the name of the stop
	Name string `json:"name"`

	// Location is the geographical location of the stop
	Location struct {
		// Lat is the latitude of the stop
		Lat float64 `json:"lat"`
		// Lon is the longitude of the stop
		Lon float64 `json:"lon"`
	} `json:"location"`
}
