package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	// "github.com/satori/go.uuid"
	"github.com/stripe/stripe-go"
	// customer "github.com/stripe/stripe-go/customer"
	currency "github.com/stripe/stripe-go/currency"
	order "github.com/stripe/stripe-go/order"
	product "github.com/stripe/stripe-go/product"
	"net/http"
	"os"
	// "time"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func init() {
	stripeAccessToken := os.Getenv("STRIPE_KEY")
	stripe.Key = stripeAccessToken
}

func fetchAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("limit", "", "3")
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

func manageOrder(w http.ResponseWriter, r *http.Request) {
	// if prevSession, ok := session.Values["session_id"].(string); !ok || !auth {

	// 	return
	// }
	return
}

func secret(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session)

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(r)
	fmt.Println("session", session)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["session"] = true
	session.Save(r, w)

}
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	fmt.Println("session", session)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = false
	session.Save(r, w)
}

// func createSessionId() string {
// 	uid, _ := uuid.NewV4()
// 	return uid
// }

func createCustomer(w http.ResponseWriter, r *http.Request) {
	// sessionId := createSessionId()
	// customerParams := &stripe.CustomerParams{}
	// c, err := customer.New(customerParams)
	return

}

func getCustomer(lastSessionId string) {

}

func getOrder(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

}

func fetchOrderById(orderId string) stripe.Order {
	o, err := order.Get(orderId, nil)
	fmt.Println(err)
	return *o

}

func createOrder(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session)

	// Check if user is authenticated
	if prevOrder, ok := session.Values["order-id"].(string); ok {
		fmt.Println("ALREADY HAS AN ORDER, fetching order", prevOrder)
		order := fetchOrderById(prevOrder)
		js, err := json.Marshal(order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	newOrder, _ := order.New(&stripe.OrderParams{
		Currency: currency.USD,
		Items: []*stripe.OrderItemParams{
			&stripe.OrderItemParams{
				Type:   "sku",
				Parent: "WR140S",
			},
		},
		Shipping: &stripe.ShippingParams{
			Name: "Matthew Miller",
			Address: &stripe.AddressParams{
				Line1:      "1234 Main Street",
				City:       "San Francisco",
				State:      "CA",
				Country:    "US",
				PostalCode: "94111",
			},
		},
		Email: "matthew.miller@example.com",
	})

	session.Values["order-id"] = newOrder.ID
	session.Save(r, w)

	js, err := json.Marshal(newOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
