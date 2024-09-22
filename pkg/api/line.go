package api

type Line struct {
	// ID is the unique identifier of the line
	ID int `json:"id"`
	// Name is the name of the line provided by the bus company
	Name string `json:"name"`
}
