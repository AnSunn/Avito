package responses

type Operation_e string

const (
	Add    = "Active"
	Delete = "Removed"
)

type DataByPeriod struct {
	User_id       int
	Segment_title string
	Operation     Operation_e
	Date          string
}
