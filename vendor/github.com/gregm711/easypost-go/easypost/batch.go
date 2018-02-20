package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#batches
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Batch struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Reference string    `json:"reference"`
	State     string    `json:"state"`

	NumShipments int             `json:"num_shipments"`
	Shipments    []BatchShipment `json:"shipments"`
	LabelUrl     string          `json:"label_url"`
	ScanForm     *ScanForm       `json:"scan_form"`
	Pickup       Pickup          `json:"pickup"`
	Error        *Error          `json:"error"`
}

type BatchShipment struct {
	ID           string `json:"id"` //Shipment ID
	Reference    string `json:"reference"`
	BatchStatus  string `json:"batch_status"`
	BatchMessage string `json:"batch_message"`
}

//AddShipment is adding one or multiple shipments to a batch
func (b *Batch) AddShipment(shipments []string) error {
	bodyString := ""
	for i := range shipments {
		bodyString = fmt.Sprintf("%v&batch[shipments][%v][id]=%v", bodyString, i, shipments[i])
	}
	obj, err := Request.do("POST", "shipment", fmt.Sprintf("%v/add_shipments", b.ID), bodyString)
	if err != nil {
		return errors.New("Failed to request EasyPost shipment insurance")
	}
	return json.Unmarshal(obj, &b)
}

//Create is creating an EasyPost batch
func (b *Batch) Create() error {
	obj, err := Request.do("POST", "batch", "", b.getCreatePayload("batch"))
	if err != nil {
		return errors.New("Failed to request EasyPost batch creation")
	}
	return json.Unmarshal(obj, &b)
}

func (b *Batch) GenerateScanForm() error {
	obj, err := Request.do("POST", "batch", fmt.Sprintf("%v/scan_form", b.ID), "")
	if err != nil {
		return errors.New("Failed to request EasyPost shipment insurance")
	}
	return json.Unmarshal(obj, &b)
}

func (b *Batch) Get() error {
	obj, err := Request.do("GET", "batch", b.ID, "")
	if err != nil {
		return errors.New("Failed to retrieve EasyPost batch")
	}
	return json.Unmarshal(obj, &b)
}

func (b Batch) getCreatePayload(prefix string) string {
	bodyString := ""

	for i := range b.Shipments {
		bodyString = fmt.Sprintf("%v&%v[shipments][%v][id]=%v", bodyString, prefix, i, b.Shipments[i].ID)
	}

	return bodyString
}
