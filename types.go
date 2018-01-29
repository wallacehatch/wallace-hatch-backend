package main

type accountRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type cardRequest struct {
	Number string `json:"cardNumber"`
	Exp    string `json:"exp"`
	CVC    string `json:"cvc"`
}

type orderRequest struct {
	Items []itemRequest `json:"items"`
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

type completeOrderRequest struct {
	Account  accountRequest  `json:"account"`
	Card     cardRequest     `json:"cardInfo"`
	Order    orderRequest    `json:"cart"`
	Shipping shippingRequest `json:"shipping"`
	Coupon   string          `json:"coupon"`
}
