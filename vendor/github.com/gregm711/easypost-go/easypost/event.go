package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#events
*/

import "time"

const (
	webHookObjectTracker = "Tracker"
)

//Event is created by changes in objects created via the API
type Event struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Description        string     `json:"description"`
	PreviousAttributes Attributes `json:"previous_attributes"`
	PendingUrls        []string   `json:"pending_urls"`
	CompletedUrls      []string   `json:"completed_urls"`
	Result             Tracker    `json:"result"`
}

//Attributes are attributes
type Attributes struct {
	Status string `json:"status"`
}

//NewEvent returns a new istance of Event
func NewEvent(id string, createdAt, updatedAt time.Time) Event {
	return Event{
		ID:        id,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
		Object:    "Event",
	}
}
