package main

type accountRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	AcceptTerms bool   `json:"acceptTerms"`
}

type cardRequest struct {
	Number string `json:"cardNumber"`
	Exp    string `json:"exp"`
	CVC    string `json:"cvc"`
}

type orderRequest struct {
	Items []itemRequest `json:"items"`
}

type couponSubmitRequest struct {
	Email string `json:"email"`
}

type itemRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type shippingRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	StreetNumber string `json:"streetNumber"`
	StreetName   string `json:"streetName"`
	AptSuite     string `json:"aptSuite"`
	Company      string `json:"company"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Zip          string `json:"zip"`
}

type Response struct {
	Text   string `json:"text"`
	Status int    `json:"status"`
}

type completeOrderRequest struct {
	Account     accountRequest         `json:"account"`
	Card        cardRequest            `json:"cardInfo"`
	Order       orderRequest           `json:"cart"`
	Shipping    shippingRequest        `json:"shipping"`
	Coupon      string                 `json:"coupon"`
	GooglePlace map[string]interface{} `json:"google_place"`
}

type Email struct {
	From      string `json:"from" bson:"from"`
	To        string `json:"to" bson:"to"`
	Subject   string `json:"subject" bson:"subject"`
	Html      string `json:"html" bson:"html"`
	PlainText string `json:"plain_text" bson:"plain_text"`
}

type EmailItemInformation struct {
	Name     string
	Size     string
	Color    string
	Style    string
	ImageUrl string
	Quantity int
	Price    float64
}

type EmailShippingInformation struct {
	Address          string
	City             string
	State            string
	Zip              string
	TrackingNumber   string
	TrackingCarrier  string
	EstimatedArrival string // dont need to get fancy ,just set to 4-7 days for now
	TrackingUrl      string
}

type EmailInformation struct {
	To             string
	From           string
	FirstName      string
	OrderNumber    string
	OrderDate      string
	OrderTotal     float64
	Items          []EmailItemInformation
	CardType       string
	CardImageUrl   string
	CardMask       string
	Shipping       EmailShippingInformation
	NumItems       int
	NumItemsMinus  int
	CouponCode     string
	CouponDiscount int
}

type ResponseError struct {
	ErrorMsg string `json:"error_message"`
	Status   int    `json:"status"`
}

type idsReqeust struct {
	Ids []string `json:"product_ids"`
}

// // Event is the resource representing a Stripe event.
// // For more details see https://stripe.com/docs/api#events.
// type Event struct {
// 	Account  string        `json:"account"`
// 	Created  int64         `json:"created"`
// 	Data     *EventData    `json:"data"`
// 	ID       string        `json:"id"`
// 	Live     bool          `json:"livemode"`
// 	Request  *EventRequest `json:"request"`
// 	Type     string        `json:"type"`
// 	Webhooks uint64        `json:"pending_webhooks"`
// }

// // EventRequest contains information on a request that created an event.
// type EventRequest struct {
// 	// ID is the request ID of the request that created an event, if the event
// 	// was created by a request.
// 	ID string `json:"id"`

// 	// IdempotencyKey is the idempotency key of the request that created an
// 	// event, if the event was created by a request and if an idempotency key
// 	// was specified for that request.
// 	IdempotencyKey string `json:"idempotency_key"`
// }

// // EventData is the unmarshalled object as a map.
// type EventData struct {
// 	Obj  map[string]interface{}
// 	Prev map[string]interface{} `json:"previous_attributes"`
// 	Raw  json.RawMessage        `json:"object"`
// }

// // EventList is a list of events as retrieved from a list endpoint.
// type EventList struct {
// 	ListMeta
// 	Values []*Event `json:"data"`
// }
