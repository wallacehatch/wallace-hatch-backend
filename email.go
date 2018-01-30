package main

import (
	"bytes"
	"fmt"
	"github.com/stripe/stripe-go"
	"gopkg.in/mailgun/mailgun-go.v1"
	// "html/template"
	"os"
	"strings"
	"time"
)

var mailgunApiSecretKey string
var mailgunPublicKey string
var domain string
var emailSender string

var b bytes.Buffer

func init() {
	mailgunApiSecretKey = os.Getenv("MAILGUN_API_SECRET_KEY")
	mailgunPublicKey = os.Getenv("MAILGUN_API_PUBLIC_KEY")
	domain = "mg.wallacehatch.com"
	emailSender = "info@wallacehatch.com"
}

func nameParser(fullName string) (firsName string, lastName string) {
	splits := strings.Split(fullName, " ")
	return splits[0], splits[len(splits)-1]
}

func constructEmailInformation(event stripe.Event) EmailInformation {
	emailInfo := EmailInformation{}
	orderObject := event.Data.Obj
	createdAt := int64(orderObject["created"].(float64))
	dateString := time.Unix(createdAt, 0).Format("2006-01-02")
	dateStringFormatted := fmt.Sprint(strings.Split(dateString, "-")[1], "/", strings.Split(dateString, "-")[2], "/", strings.Split(dateString, "-")[0])

	customerId, ok := orderObject["customer"].(string)
	if !ok {
		fmt.Println("no customer id ")
		customerId = "cus_CE6XPDfGQouohF"
	}
	stripeCustomer, _ := fetchCustomerFromId(customerId)
	card, _ := fetchCard(customerId, stripeCustomer.DefaultSource.ID)

	emailInfo.To = stripeCustomer.Email
	emailInfo.CardMask = card.LastFour
	emailInfo.CardType = string(card.Brand)
	emailInfo.OrderDate = dateStringFormatted
	shippingInfo := orderObject["shipping"].(map[string]interface{})
	addressInfo := shippingInfo["address"].(map[string]interface{})

	emailInfo.Shipping.Address = addressInfo["line1"].(string)
	emailInfo.Shipping.City = addressInfo["city"].(string)
	emailInfo.Shipping.State = addressInfo["state"].(string)
	emailInfo.Shipping.Zip = addressInfo["postal_code"].(string)
	emailInfo.Shipping.EstimatedArrival = "4-7"

	items := orderObject["items"].([]interface{})
	emailItems := make([]EmailItemInformation, 0)
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		if itemMap["object"].(string) == "order_item" {
			emailItem := EmailItemInformation{}
			productId, ok := itemMap["parent"].(string)
			if ok {
				skuInfo, _ := fetchProductBySku(productId)
				emailItem.ImageUrl = skuInfo.Image
				fmt.Println(skuInfo)
			}

			emailItem.Name = itemMap["description"].(string)
			quantity, ok := itemMap["quantity"].(float64)
			if !ok {
				quantity = 0
			}
			emailItem.Quantity = int(quantity)
			emailItems = append(emailItems, emailItem)
		}

	}

	emailInfo.Items = emailItems

	stripeCustomerName, ok := orderObject["name"].(string)
	if !ok {
		fmt.Println("no customer name")
		stripeCustomerName = "Greg Miller"
	}

	firstName, _ := nameParser(stripeCustomerName)
	emailInfo.FirstName = firstName
	emailInfo.OrderNumber = orderObject["id"].(string)
	emailInfo.OrderTotal = orderObject["amount"].(float64)

	return emailInfo

}

func MailgunSendEmail(email Email) (res string, err error) {
	mg := mailgun.NewMailgun(domain, mailgunApiSecretKey, mailgunPublicKey)
	message := mailgun.NewMessage(
		email.From,
		email.Subject,
		email.PlainText,
		email.To)
	message.SetTracking(true)
	message.AddBCC("greg@wallacehatch.com")
	message.SetTrackingClicks(true)
	message.SetTrackingOpens(true)
	if email.Html != "" {
		fmt.Println("setting html", email.Html)
		message.SetHtml(email.Html)
	}
	_, id, err := mg.Send(message)
	fmt.Println(id, err)
	return id, err
}

// "{
//   "created": 1326853478,
//   "livemode": false,
//   "id": "evt_00000000000000",
//   "type": "order.updated",
//   "object": "event",
//   "request": null,
//   "pending_webhooks": 1,
//   "api_version": "2018-01-23",
//   "data": {
//     "object": {
//       "id": "or_00000000000000",
//       "object": "order",
//       "amount": 13999,
//       "amount_returned": null,
//       "application": null,
//       "application_fee": null,
//       "charge": null,
//       "created": 1516897152,
//       "currency": "usd",
//       "customer": null,
//       "email": "greg711miller@gmail.com",
//       "items": [
//         {
//           "object": "order_item",
//           "amount": 13999,
//           "currency": "usd",
//           "description": "Kalio Rose",
//           "parent": "WR140S",
//           "quantity": 1,
//           "type": "sku"
//         },
//         {
//           "object": "order_item",
//           "amount": 0,
//           "currency": "usd",
//           "description": "Taxes (included)",
//           "parent": null,
//           "quantity": null,
//           "type": "tax"
//         },
//         {
//           "object": "order_item",
//           "amount": 0,
//           "currency": "usd",
//           "description": "Free shipping",
//           "parent": "ship_free-shipping",
//           "quantity": null,
//           "type": "shipping"
//         }
//       ],
//       "livemode": false,
//       "metadata": {
//       },
//       "returns": {
//         "object": "list",
//         "data": [

//         ],
//         "has_more": false,
//         "total_count": 0,
//         "url": "/v1/order_returns?order=or_1BoDRUGPb2UAQvII5Reaf1Ys"
//       },
//       "selected_shipping_method": "ship_free-shipping",
//       "shipping": {
//         "address": {
//           "city": "San Francisco",
//           "country": "US",
//           "line1": "1234 Main Street",
//           "line2": null,
//           "postal_code": "94111",
//           "state": "CA"
//         },
//         "carrier": null,
//         "name": "Matthew Miller",
//         "phone": null,
//         "tracking_number": null
//       },
//       "shipping_methods": [
//         {
//           "id": "ship_free-shipping",
//           "amount": 0,
//           "currency": "usd",
//           "delivery_estimate": null,
//           "description": "Free shipping"
//         }
//       ],
//       "status": "canceled",
//       "status_transitions": {
//         "canceled": 1516910957,
//         "fulfiled": null,
//         "paid": null,
//         "returned": null
//       },
//       "updated": 1516910957
//     },
//     "previous_attributes": {
//     }
//   }
// }"
