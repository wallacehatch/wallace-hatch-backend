package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	// log "github.com/sirupsen/logrus".
	"github.com/stripe/stripe-go"
	card "github.com/stripe/stripe-go/card"
	coupon "github.com/stripe/stripe-go/coupon"
	currency "github.com/stripe/stripe-go/currency"
	customer "github.com/stripe/stripe-go/customer"
	order "github.com/stripe/stripe-go/order"
	product "github.com/stripe/stripe-go/product"
	"net/http"
	"os"
	"strings"
	// "time"
)

func init() {
	stripeAccessToken := os.Getenv("STRIPE_KEY")
	stripe.Key = stripeAccessToken
}

func fetchAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	params := &stripe.ProductListParams{}
	products := make([]*stripe.Product, 0)
	i := product.List(params)
	for i.Next() {
		p := i.Product()
		products = append(products, p)

	}
	js, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func fetchProductByIdHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	p, err := product.Get(vars["key"], nil)
	js, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func fetchOrCreateCustomer(email string) *stripe.Customer {
	currentCustomer := &stripe.Customer{}
	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("email", "", email)
	i := customer.List(params)
	for i.Next() {
		currentCustomer = i.Customer()
	}
	if currentCustomer.ID == "" {
		customerParams := &stripe.CustomerParams{
			Email: email,
		}
		currentCustomer, _ = customer.New(customerParams)
	}
	return currentCustomer

}

func getOrder(orderId string) (*stripe.Order, error) {
	o, err := order.Get(orderId, nil)
	return o, err

}

type idsReqeust struct {
	Ids []string `json:"product_ids"`
}

func fetchProductsByIds(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t idsReqeust
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	params := &stripe.ProductListParams{}
	products := make([]*stripe.Product, 0)
	i := product.List(params)

	for i.Next() {
		for _, id := range t.Ids {
			if i.Product().ID == id {
				products = append(products, i.Product())
			}
		}
	}
	js, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func payOrder(orderId string, customerId string) (*stripe.Order, error) {
	orderPayParams := &stripe.OrderPayParams{}
	orderPayParams.Customer = customerId
	o, err := order.Pay(
		orderId,
		orderPayParams,
	)
	return o, err
}

func mapToStripeOrderParams(currentOrder orderRequest, currentShipping shippingRequest, customerId string, coupon string) *stripe.OrderParams {
	mappedOrder := &stripe.OrderParams{}
	mappedOrder.Currency = currency.USD
	orderItemParams := make([]*stripe.OrderItemParams, 0)

	address := &stripe.AddressParams{
		Line1:      currentShipping.Address,
		City:       currentShipping.City,
		State:      currentShipping.State,
		Country:    "US",
		PostalCode: currentShipping.Zip,
	}

	for _, item := range currentOrder.Items {
		quant := int64(item.Quantity)
		stripeItem := &stripe.OrderItemParams{Type: "sku", Parent: item.SKU, Quantity: &quant}
		orderItemParams = append(orderItemParams, stripeItem)

	}
	mappedOrder.Items = orderItemParams
	mappedOrder.Shipping = &stripe.ShippingParams{Name: currentShipping.Name, Address: address}
	mappedOrder.Customer = customerId
	mappedOrder.Coupon = coupon

	return mappedOrder
}

func mapToStripeCustomer(account accountRequest) stripe.Customer {
	mappedCustomer := stripe.Customer{}
	mappedCustomer.Email = account.Email
	mappedCustomer.Meta = map[string]string{"name": account.Name}
	return mappedCustomer
}

func mapToStripeCardParams(cardInfo cardRequest, customer stripe.Customer) stripe.CardParams {
	mappedCard := stripe.CardParams{}
	mappedCard.Number = cardInfo.Number
	dateSplit := strings.Split(cardInfo.Exp, "/")
	year := fmt.Sprint("20", dateSplit[1])
	mappedCard.Month = dateSplit[0]
	mappedCard.Year = year
	mappedCard.CVC = cardInfo.CVC
	mappedCard.Customer = customer.ID
	return mappedCard
}

func fetchCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := coupon.Get(vars["key"], nil)
	js, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

// create customer, create card, create order, pay order
func submitOrder(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var orderRequest completeOrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	customer := fetchOrCreateCustomer(orderRequest.Account.Email)
	stripeCardParams := mapToStripeCardParams(orderRequest.Card, *customer)
	_, err = card.New(&stripeCardParams)
	stripeOrderParams := mapToStripeOrderParams(orderRequest.Order, orderRequest.Shipping, customer.ID, orderRequest.Coupon)
	newOrder, err := order.New(stripeOrderParams)
	success, err := payOrder(newOrder.ID, customer.ID)
	fmt.Println(success, err)

}
