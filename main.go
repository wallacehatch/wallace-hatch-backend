package main

import (
	"encoding/json"
	"fmt"
	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// The person Type (more like an object)
type ContactForm struct {
	FullName     string `json:"full_name"`
	EmailAddress string `json:"email_address"`
	CompanyName  string `json:"company_name"`
	Message      string `json:"message"`
}

type EmailSignup struct {
	EmailAddress string `json:"email_address"`
	Subscribed   bool   `json:"subscribed"`
}

var mailchimpAPIKey string

func init() {
	mailchimpAPIKey = os.Getenv("MAILCHIMP_API")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/contact-form", ContactFormHandler).Methods("POST")
	router.HandleFunc("/email-signup", EmailSignupHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func ContactFormHandler(w http.ResponseWriter, r *http.Request) {
	err := mailchimp.SetKey(mailchimpAPIKey)
	contactForm := ContactForm{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	err = json.Unmarshal(body, &contactForm)
	if err != nil {

		fmt.Println("Error decoding contact form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	// Set request parameters.

	params := &members.NewParams{
		EmailAddress: contactForm.EmailAddress,
		Status:       members.StatusSubscribed,
	}
	// Add member to list 123456.
	member, err := members.New("123456", params)
	fmt.Println("new member is ", member)
	respondJson("true", http.StatusOK, w)
	return

}

func EmailSignupHandler(w http.ResponseWriter, r *http.Request) {
	err := mailchimp.SetKey(mailchimpAPIKey)
	emailSignup := EmailSignup{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	err = json.Unmarshal(body, &emailSignup)
	if err != nil {
		fmt.Println("Error decoding email singup form ", err)
		respondJson("false", http.StatusInternalServerError, w)
		return
	}
	// Set request parameters.
	params := &members.NewParams{
		EmailAddress: emailSignup.EmailAddress,
		Status:       members.StatusSubscribed,
	}
	// Add member to list 123456.
	member, err := members.New("35767", params)

	respondJson("true", http.StatusOK, w)
	fmt.Println("new member is ", member)
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
