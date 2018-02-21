package easypost

/*
Official API documentation available at:
https://www.easypost.com/docs/api.html#shipments
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//these are all the possible label formats
const (
	LabelFormatEPL2 = "EPL2"
	LabelFormatPDF  = "PDF"
	LabelFormatPNG  = "PNG"
	LabelFormatZPL  = "ZPL"
)

//Shipment is an EasyPost object that defines a shipment
type Shipment struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Reference     string       `json:"reference"`
	ToAddress     Address      `json:"to_address"`
	FromAddress   Address      `json:"from_address"`
	ReturnAddress Address      `json:"return_address"`
	BuyerAddress  Address      `json:"buyer_address"`
	Parcel        Parcel       `json:"parcel"`
	CustomsInfo   CustomsInfo  `json:"customs_info"`
	ScanForm      ScanForm     `json:"scan_form"`
	Forms         []ScanForm   `json:"forms"`
	Insurance     string       `json:"insurance"`
	Rates         []Rate       `json:"rates"`
	SelectedRate  Rate         `json:"selected_rate"`
	PostageLabel  PostageLabel `json:"postage_label"`
	Messages      []Message    `json:"messages"`
	Options       *Options     `json:"options"`
	IsReturn      bool         `json:"is_return"`
	TrackingCode  string       `json:"tracking_code"`
	UspsZone      int          `json:"usps_zone"`
	Status        string       `json:"status"`
	Tracker       Tracker      `json:"tracker"`
	Fees          []Fee        `json:"fees"`
	RefundStatus  string       `json:"refund_status"`
	BatchID       string       `json:"batch_id"`
	BatchStatus   string       `json:"batch_status"`
	BatchMessage  string       `json:"batch_message"`

	Error *Error `json:"error"`
}

//Message is an EasyPost object that defines the message for a shipment
type Message struct {
	Carrier string `json:"carrier"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

//ScanForm is and EasyPost object and defines a form for a shipment
type ScanForm struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Status        string   `json:"status"`
	Message       string   `json:"message"`
	Address       Address  `json:"address"`
	TrackingCodes []string `json:"tracking_codes"`
	FormURL       string   `json:"form_url"`
	FormFileType  string   `json:"form_file_type"`
	BatchID       string   `json:"batch_id"`
}

type PostageLabel struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IntegratedForm  string     `json:"integrated_form"`
	LabelDate       *time.Time `json:"label_date"`
	LabelEp12Url    string     `json:"label_ep12_url"`
	LabelFileType   string     `json:"label_file_type"`
	LabelPdfURL     string     `json:"label_pdf_url"`
	LabelResolution int        `json:"label_resolution"`
	LabelSize       string     `json:"label_size"`
	LabelType       string     `json:"label_type"`
	LabelURL        string     `json:"label_url"`
	LabelZplURL     string     `json:"label_zpl_url"`
}

//Create an EasyPost shipment
func (s *Shipment) Create() error {
	obj, err := Request.do("POST", "shipment", "", s.getCreatePayload("shipment"))
	if err != nil {
		return errors.New("Failed to request EasyPost shipment creation")
	}
	return json.Unmarshal(obj, &s)
}

//Buy an EasyPost shipment
func (s *Shipment) Buy() error {
	rate := s.SelectedRate.Service
	carrier := s.SelectedRate.Carrier
	if rate == "" {
		return errors.New("No rate has been selected")
	}
	if s.ID == "" {
		if err := s.Create(); err != nil {
			return err
		}
	}

	if s.SelectedRate.ID == "" {

		for _, singleRate := range s.Rates {
			if rate == singleRate.Service && carrier == singleRate.Carrier {
				s.SelectedRate = singleRate
				break
			}
		}
		if s.SelectedRate.ID == "" {
			return errors.New("Cannot find rate '" + rate + "' for the given carrier '" + carrier + "'")
		}

	}

	obj, err := Request.do("POST", "shipment", fmt.Sprintf("%v/buy", s.ID), fmt.Sprintf("rate[id]=%v", s.SelectedRate.ID))
	if err != nil {
		return errors.New("Failed to request EasyPost shipment purchase")
	}

	err = json.Unmarshal(obj, &s)
	if err != nil {
		return errors.New("Failed to decode EasyPost shipment purchase")
	}

	if s.Error != nil {
		return fmt.Errorf("Failed to request EasyPost shipment purcahse: %v", s.Error.Message)
	}

	return nil
}

//ConvertLabel is requesting shipping label conversion to the given file format
func (s *Shipment) ConvertLabel(format string) error {
	obj, err := Request.do("GET", "shipment", fmt.Sprintf("%v/label?file_format=%v", s.ID, format), "")
	if err != nil {
		return fmt.Errorf("Failed to request EasyPost shipping label conversion : %v", err)
	}
	return json.Unmarshal(obj, &s)
}

func (s *Shipment) Get() error {
	obj, err := Request.do("GET", "shipment", fmt.Sprintf(s.ID), "")
	if err != nil {
		return fmt.Errorf("Failed to fetch shipment from ID: %", err)
	}
	return json.Unmarshal(obj, &s)
}

//Insure is requesting insurance for the given amount
func (s *Shipment) Insure(amount float32) error {
	obj, err := Request.do("POST", "shipment", fmt.Sprintf("%v/insure", s.ID), fmt.Sprintf("amount=%v", amount))
	if err != nil {
		return errors.New("Failed to request EasyPost shipment insurance")
	}
	return json.Unmarshal(obj, &s)
}

//getCreatePayload returns the payload to append to the EasyPost API request
func (s Shipment) getCreatePayload(prefix string) string {
	bodyString := ""

	bodyString = fmt.Sprintf("%v&%vfile_format=%v", bodyString, prefix, LabelFormatPNG)
	bodyString = fmt.Sprintf("%v&%v[from_address][id]=%v", bodyString, prefix, s.FromAddress.ID)
	bodyString = fmt.Sprintf("%v&%v[to_address][id]=%v", bodyString, prefix, s.ToAddress.ID)
	var parcelPrefix = fmt.Sprintf("%v[parcel]", prefix)
	bodyString = fmt.Sprintf("%v&%v", bodyString, s.Parcel.getCreatePayload(parcelPrefix))
	bodyString = fmt.Sprintf("%v&%v[carrier_accounts][0]=%v", bodyString, prefix, s.SelectedRate.CarrierAccountID)
	bodyString = fmt.Sprintf("%v&%v[customs_info][id]=%v", bodyString, prefix, s.CustomsInfo.ID)
	if s.Reference != "" {
		bodyString = fmt.Sprintf("%v&%v[reference]=%v", bodyString, prefix, s.Reference)
	}

	if s.Options != nil {
		var optionsPrefix = fmt.Sprintf("%v[options]", prefix)
		bodyString = fmt.Sprintf("%v&%v", bodyString, s.Options.getCreatePayload(optionsPrefix))
	}

	if s.IsReturn {
		bodyString = fmt.Sprintf("%v&%v[is_return]=true", bodyString, prefix)
	}

	return bodyString
}
