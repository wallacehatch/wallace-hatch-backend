package main

import (
	"bytes"
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var stripeWebhookSig string

var b bytes.Buffer

func init() {

	stripeWebhookSig = os.Getenv("STRIPE_WEBHOOK_SIG")
}

func StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling webhook")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), stripeWebhookSig)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "order.created":
		fmt.Println("order was created")
	case "order.updated":
		orderUpdatedEvent(event)
	case "order.payment_succeeded":
		orderConfirmationEmail(event)
	default:
		fmt.Println("default")

	}

}

func orderUpdatedEvent(event stripe.Event) {
	statusTransitions, ok := event.Data.Prev["status_transitions"].(map[string]interface{})
	if !ok {
		fmt.Println("order has no status tranditions")
	}
	fulfiledChange, ok := statusTransitions["fulfiled"]
	if ok {
		fmt.Println("orders fulfillment updated!", fulfiledChange)
		orderShippedEmail(event)
	}
}

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}

func orderConfirmationEmail(event stripe.Event) {

	bufferBytes := bytes.Buffer{}
	emailInfo := constructEmailInformation(event)
	emailInfo.NumItemsMinus = emailInfo.NumItems - 1
	tmpl, err := template.ParseFiles("email-templates/order-confirmation.html")
	if err != nil {
		logger.Error("error opening template ", err)

	}
	if err := tmpl.Execute(&bufferBytes, emailInfo); err != nil {
		logger.Error("error executing html ", err)

	}

	email := Email{}
	email.Subject = "Order Confirmation"
	email.From = "info@wallacehatch.com"
	email.To = emailInfo.To
	email.Html = bufferBytes.String()

	MailgunSendEmail(email)

}

func orderShippedEmail(event stripe.Event) {

	constructEmailInformation(event)
	bufferBytes := bytes.Buffer{}
	emailInfo := constructEmailInformation(event)
	emailInfo.NumItemsMinus = emailInfo.NumItems - 1
	tmpl, err := template.ParseFiles("email-templates/order-shipped.html")
	if err != nil {
		logger.Error("error opening template ", err)

	}
	if err := tmpl.Execute(&bufferBytes, emailInfo); err != nil {
		logger.Error("error executing html ", err)

	}

	email := Email{}
	email.Subject = "Order Shipped"
	email.From = "info@wallacehatch.com"
	email.To = emailInfo.To
	email.Html = bufferBytes.String()
	MailgunSendEmail(email)

}
