package api

type Schedule struct {

	// Line is the line that the schedule is for
	Line Line `json:"line"`

	// Route is the route that the schedule is for
	Route string `json:"route"`

	// Time is the time of the schedule
	Time int `json:"time"`
}
