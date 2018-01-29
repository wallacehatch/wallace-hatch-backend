package main

import (
	"encoding/json"
	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

// The person Type (more like an object)
type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Message string `json:"message"`
}

type EmailSignup struct {
	Email      string `json:"email"`
	Subscribed bool   `json:"subscribed"`
}

var mailchimpAPIKey string

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	mailchimpAPIKey = os.Getenv("MAILCHIMP_API")
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	router := mux.NewRouter()
	router.HandleFunc("/contact-form/", ContactFormHandler).Methods("POST")
	router.HandleFunc("/email-signup/", EmailSignupHandler).Methods("POST")
	router.HandleFunc("/health-check/", HealthCheckHandler).Methods("GET")
	router.HandleFunc("/get-all-products/", fetchAllProductsHandler).Methods("GET")
	router.HandleFunc("/get-products/", fetchProductsByIds).Methods("POST")
	router.HandleFunc("/get-product/{key}", fetchProductByIdHandler).Methods("GET")
	router.HandleFunc("/submit-order/", submitOrder).Methods("POST")
	router.HandleFunc("/get-coupon/{key}", fetchCoupon)
	router.HandleFunc("/stripe-webhook", StripeWebhookHandler)
	handler := c.Handler(router)
	port := ":8081"
	log.Info("Serving on ", port)
	log.Fatal(http.ListenAndServe(port, handler))

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health check"))
	return
}

func ContactFormHandler(w http.ResponseWriter, r *http.Request) {
	err := mailchimp.SetKey(mailchimpAPIKey)
	contactForm := ContactForm{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Error in body form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	err = json.Unmarshal(body, &contactForm)
	if err != nil {
		log.Error("Error decoding contact form form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	status := "unsubscribed"
	mergeFields := make(map[string]interface{})
	mergeFields["COMPANY"] = contactForm.Company
	mergeFields["MESSAGE"] = contactForm.Message

	params := &members.NewParams{}
	params.EmailAddress = contactForm.Email
	params.MergeFields = mergeFields
	params.Status = members.Status(status)

	member, err := members.New("7343633629", params)
	if err != nil {
		log.Error("Error with mailchimp on contact form", err, params, member)
	}
	respondJson("true", http.StatusOK, w)
	return

}

func EmailSignupHandler(w http.ResponseWriter, r *http.Request) {
	err := mailchimp.SetKey(mailchimpAPIKey)
	emailSignup := EmailSignup{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Error in body  of email signup form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	err = json.Unmarshal(body, &emailSignup)
	if err != nil {
		log.Error("Error decoding email singup form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	status := "unsubscribed"
	if emailSignup.Subscribed {
		status = "subscribed"
	}
	params := &members.NewParams{}
	params.EmailAddress = emailSignup.Email
	params.Status = members.Status(status)
	member, err := members.New("06e5278452", params)
	if err != nil {
		log.Error("Error with mailchimp  on email form", err, params, member)
	}
	log.Info("subscribed is ", params.Status, "emailSingup is", emailSignup)
	respondJson("true", http.StatusOK, w)
	return

}

type Response struct {
	Text   string `json:"text"`
	Status int    `json:"status"`
}

func respondJson(text string, status int, w http.ResponseWriter) {

	response := Response{Text: text, Status: status}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}
