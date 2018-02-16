package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#customs-infos
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type CustomsInfo struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EelPfc              string        `json:"eel_pfc"`
	ContentsType        string        `json:"contents_type"`
	ContentsExplanation string        `json:"contents_explanation"`
	CustomsCertify      bool          `json:"customs_certify"`
	CustomsSigner       string        `json:"customs_signer"`
	NonDeliveryOption   string        `json:"non_delivery_option"`
	RestrictionType     string        `json:"restriction_type"`
	RestrictionComments string        `json:"restriction_comments"`
	CustomsItems        []CustomsItem `json:"customs_items"`

	Error *Error `json:"error"`
}

func (ci *CustomsInfo) Create() error {
	obj, err := Request.do("POST", "customs_info", "", ci.getCreatePayload("customs_info"))
	if err != nil {
		return errors.New("Failed to request EasyPost customs info creation")
	}
	return json.Unmarshal(obj, &ci)
}

//getCreatePayload returns the payload to append to the EasyPost API request
func (ci CustomsInfo) getCreatePayload(prefix string) string {
	bodyString := ""

	bodyString = fmt.Sprintf("%v&%v[contents_type]=%v", bodyString, prefix, ci.ContentsType)
	bodyString = fmt.Sprintf("%v&%v[customs_certify]=%v", bodyString, prefix, ci.CustomsCertify)
	bodyString = fmt.Sprintf("%v&%v[customs_signer]=%v", bodyString, prefix, ci.CustomsSigner)
	bodyString = fmt.Sprintf("%v&%v[eel_pfc]=%v", bodyString, prefix, ci.EelPfc)
	bodyString = fmt.Sprintf("%v&%v[non_delivery_option]=%v", bodyString, prefix, ci.NonDeliveryOption)
	bodyString = fmt.Sprintf("%v&%v[restriction_type]=%v", bodyString, prefix, ci.RestrictionType)

	for i, item := range ci.CustomsItems {
		bodyString = fmt.Sprintf("%v&%v", bodyString, item.getCreatePayload(fmt.Sprintf("%v[customs_items][%v]", prefix, i)))
	}

	return bodyString
}
