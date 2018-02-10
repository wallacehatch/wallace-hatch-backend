package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	card "github.com/stripe/stripe-go/card"
	coupon "github.com/stripe/stripe-go/coupon"
	currency "github.com/stripe/stripe-go/currency"
	customer "github.com/stripe/stripe-go/customer"
	order "github.com/stripe/stripe-go/order"
	product "github.com/stripe/stripe-go/product"
	sku "github.com/stripe/stripe-go/sku"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	stripeAccessToken := os.Getenv("STRIPE_KEY")
	stripe.Key = stripeAccessToken
	fmt.Println("heres token ", stripeAccessToken)
}

func fetchAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting all products")
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

func fetchOrCreateCustomer(email string) (*stripe.Customer, error) {
	var err error
	currentCustomer := &stripe.Customer{}
	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("email", "", email)
	i := customer.List(params)
	for i.Next() {
		currentCustomer = i.Customer()
	}
	if currentCustomer.ID == "" {
		if email == "" {
			logger.Error("No email provided")
			return currentCustomer, errors.New("No email provided")
		}
		customerParams := &stripe.CustomerParams{
			Email: email,
		}
		currentCustomer, err = customer.New(customerParams)
		if err != nil {
			logger.Error("Error creating customer from params", err)
		}
	}
	return currentCustomer, err

}

func fetchProductById(productId string) (stripe.Product, error) {
	p, err := product.Get(productId, nil)
	if err != nil {
		logger.Error("Error fetching product from id", err)
	}
	return *p, err
}

func fetchSkuById(skuId string) (stripe.SKU, error) {
	s, err := sku.Get(skuId, nil)
	if err != nil {
		logger.Error("Error fetching sku from id", err)
	}
	return *s, err
}

func getOrder(orderId string) (*stripe.Order, error) {
	o, err := order.Get(orderId, nil)
	if err != nil {
		logger.Error("Error fetching order from id", err)
	}
	return o, err

}

type idsReqeust struct {
	Ids []string `json:"product_ids"`
}

func fetchCustomerFromId(customerId string) (stripe.Customer, error) {
	c, err := customer.Get(customerId, nil)
	if err != nil {
		logger.Error("Error fetching customer from id ", err)
	}
	return *c, err
}

func fetchCard(customerId string, cardId string) (stripe.Card, error) {
	c := &stripe.Card{}
	c, err := card.Get(
		cardId,
		&stripe.CardParams{Customer: customerId},
	)
	fmt.Println("GOT CARD!", c)
	if err != nil {
		logger.Error("Error fetching card", err)

	}
	return *c, err

}

func fetchDefaultCard(customerId string) (stripe.Card, error) {
	c, err := card.Get(
		"card_1BpiApGPb2UAQvIIulZ82rdM",
		&stripe.CardParams{Customer: "cus_CEAV0uH4PaskMF"},
	)
	if err != nil {
		logger.Error("Error fetching default card", err)
	}
	return *c, err

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
	if err != nil {
		logger.Error("Error paying order", err)
	}

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
		if quant > 0 {
			fmt.Println("quantity is ", quant)
			stripeItem := &stripe.OrderItemParams{Type: "sku", Parent: item.SKU, Quantity: &quant}
			orderItemParams = append(orderItemParams, stripeItem)
		}

	}
	mappedOrder.Items = orderItemParams
	mappedOrder.Shipping = &stripe.ShippingParams{Name: currentShipping.Name, Address: address}
	mappedOrder.Customer = customerId
	mappedOrder.Coupon = coupon
	fmt.Println("MAPPED ORDER", mappedOrder)

	return mappedOrder
}

func mapToStripeCustomer(account accountRequest) (stripe.Customer, error) {
	mappedCustomer := stripe.Customer{}
	if account.Email == "" {
		logger.Error("no email provided")
		return mappedCustomer, errors.New("No email provided")
	}
	mappedCustomer.Email = account.Email
	mappedCustomer.Meta = map[string]string{"name": account.Name, "allowTexting": strconv.FormatBool(account.AcceptTerms)}
	return mappedCustomer, nil
}

func mapToStripeCardParams(cardInfo cardRequest, customer stripe.Customer) (stripe.CardParams, error) {
	logger.Info(cardInfo)
	mappedCard := stripe.CardParams{}
	mappedCard.Number = cardInfo.Number
	dateSplit := strings.Split(cardInfo.Exp, "/")
	if (len(dateSplit) < 1) || (len(cardInfo.Exp) < 1) {
		logger.Error("date not in correct format", cardInfo.Exp)

		return mappedCard, errors.New("Date for credit card not in format MM/YY")
	}
	year := fmt.Sprint("20", dateSplit[1])
	mappedCard.Month = dateSplit[0]
	mappedCard.Year = year
	mappedCard.CVC = cardInfo.CVC
	mappedCard.Customer = customer.ID
	return mappedCard, nil
}

func fetchCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := coupon.Get(vars["key"], nil)
	if err != nil {
		logger.Error("Error decoding order request default card", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}

	js, err := json.Marshal(c)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

type ResponseError struct {
	ErrorMsg string `json:"error_message"`

	Status int `json:"status"`
}

func fetchCouponById(couponId string) (stripe.Coupon, error) {
	c, err := coupon.Get(couponId, nil)
	if err != nil {
		logger.Error("Error retreiving coupon ", err)
	}
	return *c, err

}

func respondErrorJson(err error, status int, w http.ResponseWriter) {

	response := ResponseError{ErrorMsg: err.Error(), Status: status}
	jsonResponse, jsErr := json.Marshal(response)
	if jsErr != nil {
		logger.Error("Error with json response for error message", jsErr)
		http.Error(w, jsErr.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
	return

}

func createCustomer(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var orderRequest completeOrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		logger.Error("Error decoding order request", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}

	customer, err := fetchOrCreateCustomer(orderRequest.Account.Email)
	if err != nil {
		logger.Error("Error creating customer ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	stripeCardParams, err := mapToStripeCardParams(orderRequest.Card, *customer)

	if err != nil {
		logger.Error("Error mapping to card card ", err)

		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	_, err = card.New(&stripeCardParams)
	if err != nil {
		logger.Error("Error creating card ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}

	if orderRequest.Coupon != "" {
		_, err := fetchCouponById(orderRequest.Coupon)
		if err != nil {
			logger.Error("Error with coupon", err)
			respondErrorJson(err, http.StatusBadRequest, w)
			return
		}
	}

	js, err := json.Marshal(customer)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

// create customer, create card, create order, pay order
func submitOrder(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var orderRequest completeOrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		logger.Error("Error decoding order request", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	customer, err := fetchOrCreateCustomer(orderRequest.Account.Email)
	if err != nil {
		logger.Error("Error creating customer ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	stripeCardParams, err := mapToStripeCardParams(orderRequest.Card, *customer)
	if err != nil {
		logger.Error("Error creating card ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	_, err = card.New(&stripeCardParams)
	if err != nil {
		logger.Error("Error creating card ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	stripeOrderParams := mapToStripeOrderParams(orderRequest.Order, orderRequest.Shipping, customer.ID, orderRequest.Coupon)
	newOrder, err := order.New(stripeOrderParams)
	if err != nil {
		logger.Error("Error creating order ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	order, err := payOrder(newOrder.ID, customer.ID)
	if err != nil {
		logger.Error("Error paying for order ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	js, err := json.Marshal(order)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}
