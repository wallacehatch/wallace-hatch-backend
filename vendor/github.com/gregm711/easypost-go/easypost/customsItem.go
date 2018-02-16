package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#customs-items
*/

import (
	"encoding/json"
	"errors"
	"fmt"
)

type CustomsItem struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Mode      string `json:"mode"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Description    string  `json:"description"`
	Quantity       int     `json:"quantity"`
	Value          string  `json:"value"`
	Weight         float32 `json:"weight"`
	HsTariffNumber string  `json:"hs_tariff_number"`
	OriginCountry  string  `json:"origin_country"`
	Currency       string  `json:"currency"`

	Error *Error `json:"error"`
}

func (ci *CustomsItem) Create() error {
	obj, err := Request.do("POST", "customs_item", "", ci.getCreatePayload("customs_item"))
	if err != nil {
		return errors.New("Failed to request EasyPost customs item creation")
	}
	return json.Unmarshal(obj, &ci)
}

//getCreatePayload returns the payload to append to the EasyPost API request
func (ci CustomsItem) getCreatePayload(prefix string) string {
	bodyString := ""

	bodyString = fmt.Sprintf("%v&%v[description]=%v", bodyString, prefix, ci.Description)
	bodyString = fmt.Sprintf("%v&%v[hs_tariff_number]=%v", bodyString, prefix, ci.HsTariffNumber)
	bodyString = fmt.Sprintf("%v&%v[origin_country]=%v", bodyString, prefix, ci.OriginCountry)
	bodyString = fmt.Sprintf("%v&%v[quantity]=%v", bodyString, prefix, ci.Quantity)
	bodyString = fmt.Sprintf("%v&%v[value]=%v", bodyString, prefix, ci.Value)
	bodyString = fmt.Sprintf("%v&%v[weight]=%v", bodyString, prefix, ci.Weight)

	return bodyString
}
