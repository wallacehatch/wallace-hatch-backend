package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#trackers
*/

import (
	"time"
)

//these are the possible tracker statuses
const (
	TrackerStatusUnknown            = "unknown"
	TrackerStatusPreTransit         = "pre_transit"
	TrackerStatusInTransit          = "in_transit"
	TrackerStatusOutForDelivery     = "out_for_delivery"
	TrackerStatusDelivered          = "delivered"
	TrackerStatusAvailableForPickup = "available_for_pickup"
	TrackerStatusReturnToSender     = "return_to_sender"
	TrackerStatusFailure            = "failure"
	TrackerStatusCancelled          = "cancelled"
	TrackerStatusError              = "error"
)

//CarrierUPSTrackingURL is the default tracking base URL for UPS
const CarrierUPSTrackingURL = "https://wwwapps.ups.com/tracking/tracking.cgi?tracknum="

//CarrierUSPSTrackingURL is the default tracking base URL for USPS
const CarrierUSPSTrackingURL = "https://tools.usps.com/go/TrackConfirmAction?qtc_tLabels1="

//Tracker is an EasyPost object that defines a shipping tracker
type Tracker struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TrackingCode    string            `json:"tracking_code"`
	Status          string            `json:"status"`
	NumShipments    int               `json:"num_shipments"`
	Reference       string            `json:"reference"`
	ScanForm        string            `json:"scan_form"`
	Shipments       []Shipment        `json:"shipments"`
	Pickup          string            `json:"pickup"`
	LabelURL        string            `json:"label_url"`
	SignedBy        string            `json:"signed_by"`
	Weight          float64           `json:"weight"`
	EstDeliveryDate *time.Time        `json:"est_delivery_date"`
	ShipmentID      string            `json:"shipment_id"`
	Carrier         string            `json:"carrier"`
	TrackingDetails []TrackingDetails `json:"tracking_details"`
	CarrierDetail   Carrier           `json:"carrier_detail"`
	Fees            []Fee             `json:"fees"`
}

//TrackingDetails is an EasyPost object that defines the details for a shipping tracker
type TrackingDetails struct {
	Object           string    `json:"object"`
	Message          string    `json:"message"`
	Status           string    `json:"status"`
	Datetime         time.Time `json:"datetime"`
	Source           string    `json:"source"`
	TrackingLocation Location  `json:"tracking_location"`
}

//Fee is an EasyPost object that defines the fee details for a shipping tracker
type Fee struct {
	Object   string `json:"object"`
	Type     string `json:"type"`
	Amount   string `json:"amount"`
	Charged  bool   `json:"charged"`
	Refunded bool   `json:"refunded"`
}

//Carrier is a carrier
type Carrier struct {
	Object               string     `json:"object"`
	Service              string     `json:"service"`
	ContainerType        string     `json:"container_type"`
	EstDeliveryDateLocal *time.Time `json:"est_delivery_date_local"`
	EstDeliveryTimeLocal *time.Time `json:"est_delivery_time_local"`
}

//Location is a tracker location
type Location struct {
	Object  string `json:"object"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}

//NewTracker returns a new instance of Tracker
func NewTracker(id string, createdAt, updatedAt time.Time) Tracker {
	return Tracker{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Object:    "Tracker",
	}
}
