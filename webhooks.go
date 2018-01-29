package main

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"io/ioutil"
	"net/http"
	"os"
)

var stripeWebhookSig string

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
	constructEmailInformation(event)
	// statusTransitions, ok := event.Data.Prev["status_transitions"].(map[string]interface{})
	// if !ok {
	// 	fmt.Println("order has no status tranditions")
	// }
	// fulfiledChange, ok := statusTransitions["fulfiled"]
	// if ok {
	// 	fmt.Println("orders fulfillment updated!", fulfiledChange)
	// 	orderShippedEmail(event)
	// }
}

func orderConfirmationEmail(event stripe.Event) {
	fmt.Println("order was paid for successfully, time to email!")
	constructEmailInformation(event)

}

func orderShippedEmail(event stripe.Event) {
	fmt.Println("gonna send email for shipping")
	constructEmailInformation(event)

}