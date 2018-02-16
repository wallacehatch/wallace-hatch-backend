package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#pickups
*/

import "time"

//Pickup is an easypost pickup
type Pickup struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Reference        string           `json:"reference"`
	Status           string           `json:"status"`
	MinDatetime      string           `json:"min_datetime"`
	MaxDatetime      string           `json:"max_datetime"`
	IsAccountAddress bool             `json:"is_account_address"`
	Instructions     string           `json:"instructions"`
	Messages         Message          `json:"messages"`
	Confirmation     string           `json:"confirmation"`
	Shipment         Shipment         `json:"shipment"`
	Address          Address          `json:"address"`
	CarrierAccounts  []CarrierAccount `json:"carrier_accounts"`
	PickupRates      []PickupRate     `json:"pickup_rates"`
}

//PickupRate is an easypost pickup rate
type PickupRate struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Mode      string `json:"mode"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Service  string `json:"service"`
	Carrier  string `json:"carrier"`
	Rate     string `json:"rate"`
	Currency string `json:"currency"`
	PickupID string `json:"pickup_id"`
}
