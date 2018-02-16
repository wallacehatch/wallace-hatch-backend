package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#parcels
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//Parcel is an EasyPost object that defines a shipping parcel
type Parcel struct {
	ID        string     `json:"id"`
	Object    string     `json:"object"`
	Mode      string     `json:"mode"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Length            float32 `json:"length"`
	Width             float32 `json:"width"`
	Height            float32 `json:"height"`
	PredefinedPackage string  `json:"predefined_package"`
	Weight            float32 `json:"weight"`
}

//getCreatePayload returns the payload to append to the EasyPost API request
func (p Parcel) getCreatePayload(prefix string) string {
	bodyString := ""
	if p.ID != "" {
		bodyString = fmt.Sprintf("%v&%v[id]=%v", bodyString, prefix, p.ID)
	} else {
		if p.PredefinedPackage != "" {
			bodyString = fmt.Sprintf("%v&%v[predefined_package]=%v", bodyString, prefix, p.PredefinedPackage)
		} else {
			bodyString = fmt.Sprintf("%v&%v[length]=%v", bodyString, prefix, p.Length)
			bodyString = fmt.Sprintf("%v&%v[width]=%v", bodyString, prefix, p.Width)
			bodyString = fmt.Sprintf("%v&%v[height]=%v", bodyString, prefix, p.Height)
		}
		bodyString = fmt.Sprintf("%v&%v[weight]=%v", bodyString, prefix, p.Weight)
	}
	return bodyString
}

// Create an EasyPost parcel
func (p *Parcel) Create() error {
	obj, err := Request.do("POST", "parcel", "", p.getCreatePayload("parcel"))
	if err != nil {
		return errors.New("Failed to request EasyPost parcel creation")
	}
	return json.Unmarshal(obj, &p)
}
