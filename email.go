package main

import (
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

func constructEmailInformation(event stripe.Event) (EmailInformation, error) {
	emailInfo := EmailInformation{}
	orderObject := event.Data.Obj
	createdAt := int64(orderObject["created"].(float64))
	dateString := time.Unix(createdAt, 0).Format("2006-01-02")
	dateStringFormatted := fmt.Sprint(strings.Split(dateString, "-")[1], "/", strings.Split(dateString, "-")[2], "/", strings.Split(dateString, "-")[0])

	customerId, ok := orderObject["customer"].(string)
	if !ok {
		fmt.Println("no customer id ")

		customerId = "cus_CIfUBNHrYfLsJ6"
	}

	stripeCustomer, _ := fetchCustomerFromId(customerId)
	card, err := fetchCard(customerId, stripeCustomer.DefaultSource.ID)
	if err != nil {
		logger.Error("Error with card fetching", err)
	}
	emailInfo.To = stripeCustomer.Email
	emailInfo.CardMask = card.LastFour
	emailInfo.CardType = string(card.Brand)

	switch emailInfo.CardType {
	case "Visa":
		emailInfo.CardImageUrl = "https://s3.us-east-2.amazonaws.com/wallace-hatch/visa%403x.png"
	case "Mastercard":
		emailInfo.CardImageUrl = "https://s3.us-east-2.amazonaws.com/wallace-hatch/mastercard%403x.png"
	case "Discover":
		emailInfo.CardImageUrl = "https://s3.us-east-2.amazonaws.com/wallace-hatch/discover%403x.png"
	case "American Express":
		emailInfo.CardImageUrl = "https://s3.us-east-2.amazonaws.com/wallace-hatch/amex%403x.png"
	}

	emailInfo.OrderDate = dateStringFormatted
	shippingInfo := orderObject["shipping"].(map[string]interface{})
	addressInfo := shippingInfo["address"].(map[string]interface{})

	address, ok := addressInfo["line1"].(string)
	if !ok {
		logger.Error("no address supplied from webhook")
	}
	emailInfo.Shipping.Address = address

	city, ok := addressInfo["city"].(string)
	if !ok {
		logger.Error("no city supplied from webhook")
	}
	emailInfo.Shipping.City = city

	state, ok := addressInfo["state"].(string)
	if !ok {
		logger.Error("no state supplied from webhook")
	}
	emailInfo.Shipping.State = state

	zip, ok := addressInfo["postal_code"].(string)
	if !ok {
		logger.Error("no zip supplied from webhook")
	}
	emailInfo.Shipping.Zip = zip

	emailInfo.Shipping.EstimatedArrival = "4-7"
	carrier, ok := shippingInfo["carrier"].(string)
	if !ok {
		carrier = "USPS"
	}
	tracking, ok := shippingInfo["tracking_number"].(string)
	if !ok {
		tracking = "TRACK1234"

	}
	emailInfo.Shipping.TrackingCarrier = carrier
	emailInfo.Shipping.TrackingNumber = tracking
	emailInfo.Shipping.TrackingUrl = fmt.Sprint("https://tools.usps.com/go/TrackConfirmAction?tLabels=", tracking)

	items := orderObject["items"].([]interface{})
	numItems := 0
	emailItems := make([]EmailItemInformation, 0)
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		if itemMap["type"].(string) == "sku" {
			emailItem := EmailItemInformation{}
			productId, ok := itemMap["parent"].(string)
			if ok {
				skuInfo, _ := fetchSkuById(productId)
				emailItem.Size = fmt.Sprint(skuInfo.Attrs["size"], "MM")
				emailItem.ImageUrl = skuInfo.Image
				emailItem.Price = float64(skuInfo.Price) / 100.0
				productInfo, _ := fetchProductById(skuInfo.Product.ID)
				emailItem.Color = productInfo.Meta["caseColor"]
			}

			emailItem.Style = productId

			emailItem.Name = itemMap["description"].(string)
			quantity, ok := itemMap["quantity"].(float64)
			if !ok {
				quantity = 0
			}
			emailItem.Quantity = int(quantity)
			numItems = numItems + int(quantity)
			emailItems = append(emailItems, emailItem)
		}

	}

	emailInfo.Items = emailItems
	emailInfo.NumItems = numItems

	stripeCustomerName, ok := orderObject["name"].(string)
	if !ok {
		fmt.Println("no customer name")
		stripeCustomerName = "Greg Miller"
	}

	firstName, _ := nameParser(stripeCustomerName)
	emailInfo.FirstName = firstName
	emailInfo.OrderNumber = strings.Replace(orderObject["id"].(string), "or_", "", -1)
	emailInfo.OrderTotal = orderObject["amount"].(float64) / 100.0

	return emailInfo, nil

}

func MailgunSendEmail(email Email) (res string, err error) {

	if strings.Contains(email.To, "@wallacehatch.com") {
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
		if err != nil {
			logger.Error("Error sending email from mailgun ", err)
		}
		return id, err
	}

	return "", nil
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
