package easypost

import "time"

//EasypostObject is the representation of an object base fields
type EasypostObject struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
