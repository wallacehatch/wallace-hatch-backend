package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#orders
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//Order is an Easypost order
type Order struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Reference       string           `json:"reference"`
	ToAddress       Address          `json:"to_address"`
	FromAddress     Address          `json:"from_address"`
	ReturnAddress   Address          `json:"return_address"`
	BuyerAddress    Address          `json:"buyer_address"`
	Shipments       []Shipment       `json:"shipments"`
	Rates           []Rate           `json:"rates"`
	Messages        []Message        `json:"messages"`
	IsReturn        bool             `json:"is_return"`
	CarrierAccounts []CarrierAccount `json:"carrier_accounts"`
	Options         *Options         `json:"options"`

	Carrier string `json:"carrier"`
	Service string `json:"service"`

	Error *Error `json:"error"`
}

//ErrorString returns a string representation of the error
func (o Order) ErrorString() string {
	if o.Error == nil {
		return ""
	}
	var epErr = fmt.Sprintf("%v (%v)", o.Error.Message, o.Error.Code)
	if o.Error.Code == ErrorOrderRateUnavailable {
		epErr = fmt.Sprintf("%v. Rates: ", epErr)
		if len(o.Rates) == 0 {
			epErr = fmt.Sprintf("%v None available ", epErr)
		} else {
			for _, r := range o.Rates {
				epErr = fmt.Sprintf("%v [%v] %v (%v %v).", epErr, r.Carrier, r.Service, r.Rate, r.Currency)
			}
		}
	}
	return epErr
}

//Create an EasyPost order
func (o *Order) Create() error {
	obj, err := Request.do("POST", "order", "", o.getCreatePayload())
	if err != nil {
		return errors.New("Failed to request EasyPost shipment creation")
	}
	return json.Unmarshal(obj, &o)
}

//Buy an EasyPost order
func (o *Order) Buy() error {
	if o.Carrier == "" {
		return errors.New("Carrier is missing")
	}
	if o.Service == "" {
		return errors.New("Service rate is missing")
	}

	obj, err := Request.do("POST", "order", fmt.Sprintf("%v/buy", o.ID), fmt.Sprintf("carrier=%v&service=%v", o.Carrier, o.Service))
	if err != nil {
		return errors.New("Failed to request EasyPost shipment creation")
	}
	return json.Unmarshal(obj, &o)
}

//getCreatePayload returns the payload to append to the EasyPost API request
func (o Order) getCreatePayload() string {
	bodyString := ""

	bodyString = fmt.Sprintf("%v&order[reference]=%v", bodyString, o.Reference)
	bodyString = fmt.Sprintf("%v&order[is_return]=%v", bodyString, o.IsReturn)

	if o.ToAddress.ID != "" {
		bodyString = fmt.Sprintf("%v&order[to_address][id]=%v", bodyString, o.ToAddress.ID)
	} else {
		bodyString = fmt.Sprintf("%v&%v", bodyString, o.ToAddress.getPayload("order[to_address]"))
	}
	if o.FromAddress.ID != "" {
		bodyString = fmt.Sprintf("%v&order[from_address][id]=%v", bodyString, o.FromAddress.ID)
	} else {
		bodyString = fmt.Sprintf("%v&%v", bodyString, o.FromAddress.getPayload("order[from_address]"))
	}

	for i, s := range o.Shipments {
		var shipmentPrefix = fmt.Sprintf("order[shipments][%v]", i)
		var parcelPrefix = fmt.Sprintf("%v[parcel]", shipmentPrefix)
		if s.Parcel.ID != "" {
			bodyString = fmt.Sprintf("%v&%v[id]=%v", bodyString, parcelPrefix, s.Parcel.ID)
		} else {
			bodyString = fmt.Sprintf("%v&%v", bodyString, s.Parcel.getCreatePayload(parcelPrefix))
		}
		if s.Options != nil {
			if s.Options.LabelDate != "" {
				bodyString = fmt.Sprintf("%v&%v[options][label_date]=%v", bodyString, shipmentPrefix, s.Options.LabelDate)
			}
		}

		var customsPrefix = fmt.Sprintf("order[shipments][%v][customs_info]", i)
		if s.CustomsInfo.ID != "" {
			bodyString = fmt.Sprintf("%v&%v[id]=%v", bodyString, customsPrefix, s.CustomsInfo.ID)
		} else {
			bodyString = fmt.Sprintf("%v&%v", bodyString, i, s.CustomsInfo.getCreatePayload(customsPrefix))
		}
	}

	for i, ca := range o.CarrierAccounts {
		bodyString = fmt.Sprintf("%v&order[carrier_accounts][%v][id]=%v", bodyString, i, ca.ID)
	}

	if o.Options != nil {
		if o.Options.LabelDate != "" {
			bodyString = fmt.Sprintf("%v&order[options][label_date]=%v", bodyString, o.Options.LabelDate)
		}
	}

	return bodyString
}
