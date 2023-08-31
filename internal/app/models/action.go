package models

// Action model
type Action struct {
	ID                   int `json:"id"`
	User_id              int `json:"user_id"`
	Segment_id           int
	Start_date, End_date string
	Segment_title_to_add []struct {
		Title string `json:"title"`
		Days  int    `json:"days"`
	} `json:"add_list"`
	Segment_title_to_del []string `json:"remove_list"`
}
