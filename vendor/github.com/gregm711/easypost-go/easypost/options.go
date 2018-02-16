package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#options
*/

import "fmt"

//Options is an EasyPost object and defines the list of carrier-specific options for a shipment
type Options struct {
	AdditionalHandling       bool   `json:"additional_handling"`
	AddressValidationLevel   string `json:"address_validation_level"`
	Alcohol                  bool   `json:"alchol"`
	BillReceiverAccount      string `json:"bill_receiver_account"`
	BillReceiverPostalCode   string `json:"bill_receiver_postal_code"`
	BillThirdPartyAccount    string `json:"bill_third_party_account"`
	BillThirdPartyCountry    string `json:"bill_third_party_country"`
	BillThirdPartyPostalCode string `json:"bill_third_party_postal_code"`
	ByDrone                  bool   `json:"by_drone"`
	CarbonNeutral            bool   `json:"carbon_neutral"`
	CodAmount                string `json:"cod_amount"`
	Currency                 string `json:"currency"`
	DeliveredDutyPaid        bool   `json:"delivery_duty_paid"`
	DeliveryConfirmation     string `json:"delivery_confirmation"`
	DryIce                   bool   `json:"dry_ice"`
	DryIceMedical            string `json:"dry_ice_medical"`
	DryIceWeight             string `json:"dry_icd_weight"`
	FreightCarge             int    `json:"freight_charge"`
	HandlingInstructions     string `json:"handling_instructions"`
	Hazmat                   string `json:"hazmat"`
	HoldForPickup            bool   `json:"hold_for_pickup"`
	InvoiceNumber            string `json:"invoice_number"`
	LabelDate                string `json:"label_date"`
	LabelFormat              string `json:"label_format"`
	Machinable               bool   `json:"machinable"`
	PrintCustom1             string `json:"print_custom_1"`
	PrintCustom2             string `json:"print_custom_2"`
	PrintCustom3             string `json:"print_custom_3"`
	PrintCustom1Barcode      string `json:"print_custom_1_barcode"`
	PrintCustom2Barcode      string `json:"print_custom_2_barcode"`
	PrintCustom3Barcode      string `json:"print_custom_3_barcode"`
	PrintCustom1Code         string `json:"print_custom_1_code"`
	PrintCustom2Code         string `json:"print_custom_2_code"`
	PrintCustom3Code         string `json:"print_custom_3_code"`
	SaturdayDelivery         bool   `json:"saturday_delivery"`
	SpecialRatesEligibility  string `json:"special_rates_eligibility"`
	SmartpostHub             string `json:"smartpost_hub"`
	SmartpostManifest        string `json:"smartpost_manifest"`
}

//getCreatePayload returns the payload to append to the EasyPost API request
func (o Options) getCreatePayload(prefix string) string {
	var bodyString = ""

	if o.LabelDate != "" {
		bodyString = fmt.Sprintf("%v&%v[label_date]=%v", bodyString, prefix, o.LabelDate)
	}
	if o.DeliveredDutyPaid {
		bodyString = fmt.Sprintf("%v&%v[delivered_duty_paid]=true", bodyString, prefix)
	}
	if o.DeliveryConfirmation != "" {
		bodyString = fmt.Sprintf("%v&%v[delivery_confirmation]=%v", bodyString, prefix, o.DeliveryConfirmation)
	}

	return bodyString
}
