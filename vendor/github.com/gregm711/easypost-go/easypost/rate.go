package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#rates
*/

import "time"

const (
	ErrorOrderRateUnavailable = "ORDER.RATE.UNAVAILABLE"
)

//Rate is an EasyPost object and defines the shipment rate, fetched after shipment creation
type Rate struct {
	ID        string     `json:"id"`
	Object    string     `json:"object"`
	Mode      string     `json:"mode"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Service                string     `json:"service"`
	Carrier                string     `json:"carrier"`
	CarrierAccountID       string     `json:"carrier_account_id"`
	ShipmentID             string     `json:"shipment_id"`
	Rate                   string     `json:"rate"`
	Currency               string     `json:"currency"`
	RetailRate             string     `json:"retail_rate"`
	RetailCurrency         string     `json:"retail_currency"`
	ListRate               string     `json:"list_rate"`
	ListCurrency           string     `json:"list_currency"`
	DeliveryDays           int        `json:"delivery_days"`
	DeliveryDate           *time.Time `json:"delivery_date"`
	DaliveryDateGuaranteed bool       `json:"delivery_date_guaranteed"`
}
