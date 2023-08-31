package models

// Segment model
type Segment struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool
}
