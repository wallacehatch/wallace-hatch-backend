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
var logger *log.Logger

func init() {
	mailchimpAPIKey = os.Getenv("MAILCHIMP_API")
	logger = log.New()
	logger.Formatter = new(log.JSONFormatter)

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
	router.HandleFunc("/create-customer/", createCustomer).Methods("POST")
	router.HandleFunc("/get-coupon/{key}", fetchCoupon)
	router.HandleFunc("/stripe-webhook/", StripeWebhookHandler)
	router.HandleFunc("/apply-for-coupon/", couponSignupHandler).Methods("POST")
	router.HandleFunc("/easypost-webhook/", easypostWebhookHandler)
	router.HandleFunc("/test-twilio/", testTwilioHandler)
	router.HandleFunc("/create-review/", createProductReviewHandler).Methods("POST")
	router.HandleFunc("/get-product-reviews/{key}", fetchProductReviewsHandler).Methods("GET")
	router.HandleFunc("/get-customer-orders/{key}", fetchPastOrdersHandler).Methods("GET")
	router.HandleFunc("/validate-review/", fetchPastOrdersHandler).Methods("POST")
	handler := c.Handler(router)
	port := ":8090"
	logger.Info("Serving on ", port)
	logger.Fatal(http.ListenAndServe(port, handler))

}

func testTwilioHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Testing twilio webhook")
	message := "Hey Greg - this message means twilio is working on production servers for wally 🤘"
	twilio, err := sendSMSMessage("4403966613", message)
	if err != nil {
		logger.Error("Error with twilio", err)
	}
	logger.Info("twilio info ", twilio)
	return

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
		logger.Error("Error in body form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	err = json.Unmarshal(body, &contactForm)
	if err != nil {
		logger.Error("Error decoding contact form form ", err)
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
		logger.Error("Error with mailchimp on contact form", err, params, member)
	}
	respondJson("true", http.StatusOK, w)
	return

}

func EmailSignupHandler(w http.ResponseWriter, r *http.Request) {
	err := mailchimp.SetKey(mailchimpAPIKey)
	emailSignup := EmailSignup{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("Error in body  of email signup form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	err = json.Unmarshal(body, &emailSignup)
	if err != nil {
		logger.Error("Error decoding email singup form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}

	member, err := addToMailchimpNewsletter(emailSignup.Email, "", "")

	if err != nil {
		logger.Error("Error with mailchimp  on email form", err, member)
	}

	respondJson("true", http.StatusOK, w)
	return

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

// mostRecentTrackingMessage := "Out for Delivery"
// currentLocation := "CHARLESTON SC"
// estimatedArrival := "Tuesday, February 27"
// shortenedTrackingLink, _ := shortenUrl("https://track.easypost.com/djE6dHJrX2ZjYTRjZTQyNDk2ZjQ5NjBiODkxNzQzOTQ1YWQ5OGMy")
// message := fmt.Sprint(mostRecentTrackingMessage, ": Your Wallace Hatch ⌚️📦 is on it's way!\n\nCurrent location 📍 ", currentLocation, "\n\nEstimated delivery 📅 ", estimatedArrival, ".\n\nTrack at ", shortenedTrackingLink)
// sendSMSMessage("4403966613", message)
