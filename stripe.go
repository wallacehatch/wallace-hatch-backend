package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
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
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func init() {
	stripeAccessToken := os.Getenv("STRIPE_KEY")
	stripe.Key = stripeAccessToken
}

func doesCustomerContainPastOrder(customerId string, productId string) bool {
	pastOrders, _ := fetchCustomerOrders(customerId)
	for _, order := range pastOrders {
		for _, item := range order.Items {
			if item.Quantity > 0 {
				product, _ := fetchProductFromSKU(item.Parent)
				if product.ID == productId {
					return true
				}
			}
		}
	}
	return false
}

func fetchPastOrdersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["key"]
	orders, err := fetchCustomerOrders(customerId)
	if err != nil {
		logger.Error("Error getting customer orders", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}

	js, _ := json.Marshal(orders)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

func fetchAllProductsHandler(w http.ResponseWriter, r *http.Request) {

	products := getAllProducts()
	js, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func getAllProducts() []*stripe.Product {
	params := &stripe.ProductListParams{}
	products := make([]*stripe.Product, 0)
	i := product.List(params)
	for i.Next() {
		p := i.Product()
		products = append(products, p)

	}
	return products

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

func fetchProductFromSKU(skuId string) (stripe.Product, error) {
	s, err := sku.Get(skuId, nil)
	product, err := fetchProductById(s.Product.ID)
	return product, err
}

func fetchCustomerOrders(customerId string) ([]stripe.Order, error) {
	orders := make([]stripe.Order, 0)
	params := &stripe.OrderListParams{}
	params.Filters.AddFilter("customer", "", customerId)
	i := order.List(params)
	if i.Err() != nil {
		logger.Error("Error getting customer orders ", i.Err())
	}
	for i.Next() {
		orders = append(orders, *i.Order())
	}
	return orders, i.Err()

}

func fetchOrCreateCustomer(customerParams stripe.CustomerParams) (*stripe.Customer, error) {
	var err error

	currentCustomer := &stripe.Customer{}
	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("email", "", customerParams.Email)
	i := customer.List(params)
	for i.Next() {
		currentCustomer = i.Customer()
	}
	if currentCustomer.ID == "" {
		if customerParams.Email == "" {
			logger.Error("No email provided")
			return currentCustomer, errors.New("No email provided")
		}
		currentCustomer, err = customer.New(&customerParams)
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

func getHash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	sha1 := hex.EncodeToString(h.Sum(nil))
	return sha1
}

func createCouponFromEmail(email string) (*stripe.Coupon, error) {
	hash := getHash(email)
	couponId := hash[0:6] // shorten to make  more acccesible to user
	c, err := coupon.New(&stripe.CouponParams{
		Percent:     15,
		Duration:    "once",
		ID:          couponId,
		Redemptions: 1,
		RedeemBy:    int64(time.Now().AddDate(0, 1, 0).Unix()),
	})
	return c, err
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
	if err != nil {
		logger.Error("Error fetching card", err)
	}
	return *c, err

}

func getProductsFromNames(names []string) []stripe.Product {
	products := make([]stripe.Product, 0)
	allProducts := getAllProducts()
	for _, product := range allProducts {
		cleanedProductName := cleanString(product.Name)
		for _, name := range names {
			cleanedName := cleanString(name)
			if cleanedProductName == cleanedName {
				products = append(products, *product)
			}
		}
	}
	return products
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
		logger.Error("Error fetching product by ids: ", err)
		return
	}

	products := getProductsFromIds(t.Ids)

	js, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getProductsFromIds(ids []string) []*stripe.Product {
	params := &stripe.ProductListParams{}
	products := make([]*stripe.Product, 0)
	i := product.List(params)

	for i.Next() {
		for _, id := range ids {
			if i.Product().ID == id {
				products = append(products, i.Product())
			}
		}
	}
	return products
}

func updateOrderMeta(orderId string, meta map[string]string) error {
	params := &stripe.OrderUpdateParams{}
	params.Meta = meta
	_, err := order.Update(orderId, params)
	return err
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

func mapToStripeShippingParams(googlePlace map[string]interface{}, currentShipping shippingRequest) *stripe.ShippingParams {
	mappedAddress := &stripe.AddressParams{}
	addressComponents, _ := googlePlace["address_components"].([]interface{})
	formattedAddress := googlePlace["formatted_address"].(string)
	addressSplit := strings.Split(formattedAddress, ",")

	mappedAddress.Line1 = addressSplit[0]
	mappedAddress.Line2 = currentShipping.AptSuite

	for _, addressInfo := range addressComponents {

		info := addressInfo.(map[string]interface{})
		infoType := info["types"].([]interface{})
		switch infoType[0] {
		case "country":
			mappedAddress.Country = info["short_name"].(string)
		case "postal_code":
			mappedAddress.PostalCode = info["short_name"].(string)
		case "locality":
			mappedAddress.City = info["short_name"].(string)
		case "administrative_area_level_1":
			mappedAddress.State = info["short_name"].(string)

		}
	}
	mappedShipping := &stripe.ShippingParams{Name: currentShipping.Name, Address: mappedAddress}

	return mappedShipping

}

func mapToStripeOrderParams(currentOrder orderRequest, shippingParams *stripe.ShippingParams, customerId string, coupon string) *stripe.OrderParams {
	mappedOrder := &stripe.OrderParams{}
	mappedOrder.Currency = currency.USD
	orderItemParams := make([]*stripe.OrderItemParams, 0)

	for _, item := range currentOrder.Items {

		quant := int64(item.Quantity)
		if quant > 0 {
			stripeItem := &stripe.OrderItemParams{Type: "sku", Parent: item.SKU, Quantity: &quant}
			orderItemParams = append(orderItemParams, stripeItem)
		}
	}
	mappedOrder.Items = orderItemParams
	mappedOrder.Shipping = shippingParams
	mappedOrder.Customer = customerId

	mappedOrder.Coupon = coupon

	return mappedOrder
}

func mapToStripeCustomerParams(account accountRequest) (stripe.CustomerParams, error) {
	mappedCustomer := stripe.CustomerParams{}
	if account.Email == "" {
		logger.Error("no email provided")
		return mappedCustomer, errors.New("No email provided")
	}
	mappedCustomer.Email = account.Email
	mappedCustomer.Meta = map[string]string{"name": account.Name, "allowTexting": strconv.FormatBool(account.AcceptTerms), "phone": account.Phone}
	return mappedCustomer, nil
}

func mapToStripeCardParams(cardInfo cardRequest, customer stripe.Customer) (stripe.CardParams, error) {
	mappedCard := stripe.CardParams{}
	mappedCard.Number = cardInfo.Number
	dateSplit := strings.Split(cardInfo.Exp, "/")
	if (len(dateSplit) < 1) || (len(cardInfo.Exp) < 1) {
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

func fetchOrderById(orderId string) (stripe.Order, error) {
	o, err := order.Get(orderId, nil)
	if err != nil {
		logger.Error("Error fetching stripe order by id", err)
	}
	return *o, err

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
	mappedCustomer, err := mapToStripeCustomerParams(orderRequest.Account)
	customer, err := fetchOrCreateCustomer(mappedCustomer)
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
	mappedCustomer, err := mapToStripeCustomerParams(orderRequest.Account)
	customer, err := fetchOrCreateCustomer(mappedCustomer)
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

	shippingParams := mapToStripeShippingParams(orderRequest.GooglePlace, orderRequest.Shipping)

	stripeOrderParams := mapToStripeOrderParams(orderRequest.Order, shippingParams, customer.ID, orderRequest.Coupon)

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
	go easypostController(*newOrder)

	js, err := json.Marshal(order)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

// when a customer signs up for the coupon, the are subscribed to our email list and a customer object is created,
// as well as a hashed coupon code that is emailed to be redemed
func couponSignupHandler(w http.ResponseWriter, r *http.Request) {
	bufferBytes := bytes.Buffer{}
	decoder := json.NewDecoder(r.Body)
	var couponRequest couponSubmitRequest
	err := decoder.Decode(&couponRequest)
	if err != nil {
		logger.Error("Error decoding coupon request", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	account := accountRequest{}
	account.Email = couponRequest.Email
	mappedCustomer, _ := mapToStripeCustomerParams(account)
	_, err = fetchOrCreateCustomer(mappedCustomer)
	if err != nil {
		logger.Error("Error creating customer ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	coupon, err := createCouponFromEmail(couponRequest.Email)
	if err != nil {
		logger.Error("Error creating coupon from email", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	emailInfo := EmailInformation{}
	emailInfo.To = couponRequest.Email
	emailInfo.CouponCode = coupon.ID
	emailInfo.CouponDiscount = int(coupon.Percent)
	tmpl, err := template.ParseFiles("email-templates/coupon.html")
	if err != nil {
		logger.Error("error opening template ", err)
	}
	if err := tmpl.Execute(&bufferBytes, emailInfo); err != nil {
		logger.Error("error executing html ", err)
	}
	email := Email{}
	email.Subject = fmt.Sprint("Welcome To Wallace Hatch - Take ", emailInfo.CouponDiscount, "% Off Your First Order")
	email.From = emailSender
	email.To = emailInfo.To
	email.Html = bufferBytes.String()
	_, err = MailgunSendEmail(email, applyCouponTag, time.Now())
	if err != nil {
		logger.Error("Error sending mailgun email for coupon ", err)
	}

	_, err = addToMailchimpNewsletter(couponRequest.Email, "", "")
	if err != nil {
		logger.Error("Error adding user to mailchimp newsletter ", err)
	}

	js, _ := json.Marshal(coupon)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

func cleanString(s string) string {
	cleanedString := ""
	for _, char := range s {
		if unicode.IsLetter(char) {
			lowercase := strings.ToLower(string(char))
			cleanedString = fmt.Sprint(cleanedString, lowercase)
		}
	}
	return cleanedString
}
