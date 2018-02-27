package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
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
		fmt.Println("order has no status transitions")
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
	logger.Info("order confirmation time")

	bufferBytes := bytes.Buffer{}
	emailInfo, _ := constructEmailInformationFromEvent(event)
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
	email.From = emailSender
	email.To = emailInfo.To
	email.Html = bufferBytes.String()
	MailgunSendEmail(email, orderConfirmationTag, time.Now())

}

func orderShippedEmail(event stripe.Event) {
	bufferBytes := bytes.Buffer{}
	emailInfo, _ := constructEmailInformationFromEvent(event)
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
	email.From = emailSender
	email.To = emailInfo.To
	email.Html = bufferBytes.String()
	MailgunSendEmail(email, shippedTag, time.Now())
}

func easypostWebhookHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("received webhook for easyopst:")
	var hook easypostWebhook
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&hook)
	if err != nil {
		logger.Error("Error decoding easy post webhook ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}

	shipment, _ := fetchShipmentFromId(hook.Result.ShipmentID)
	orderId := shipment.Reference

	order, err := fetchOrderById(orderId)
	if err != nil {
		logger.Error("Error fetching order from id ", err)
	}
	customer, err := fetchCustomerFromId(order.Customer.ID)
	if err != nil {
		logger.Error("Error fetching customer from ID", err)
	}
	logger.Info(customer.Meta["allowTexting"], customer.Meta["phone"])
	// customer wants to get information via sms on tracking
	if customer.Meta["allowTexting"] == "true" && customer.Meta["phone"] != "" {
		message := constructMessage(hook)
		if message != "" {
			response, err := sendSMSMessage(customer.Meta["phone"], message)
			logger.Info(response, err)
		}

	}

}
