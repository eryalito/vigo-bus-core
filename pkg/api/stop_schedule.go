package api

type StopSchedule struct {
	// Stop is the stop that the schedule is for
	Stop Stop `json:"stop"`

	// Schedules is a list of the schedules for the stop
	Schedules []Schedule `json:"schedules"`
}
