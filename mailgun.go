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

var (
	mailgunApiSecretKey string
	mailgunPublicKey    string
)

const (
	orderConfirmationTag = "orderConfirmationEmail"
	shippedTag           = "shippedEmail"
	applyCouponTag       = "applyCouponEmail"
	domain               = "mg.wallacehatch.com"
	emailSender          = "Wallace Hatch <orders@wallacehatch.com>"
)

func init() {
	mailgunApiSecretKey = os.Getenv("MAILGUN_API_SECRET_KEY")
	mailgunPublicKey = os.Getenv("MAILGUN_API_PUBLIC_KEY")
}

func nameParser(fullName string) (firsName string, lastName string) {
	splits := strings.Split(fullName, " ")
	return splits[0], splits[len(splits)-1]
}

func constructEmailInformationFromEvent(event stripe.Event) (EmailInformation, error) {
	emailInfo := EmailInformation{}
	orderObject := event.Data.Obj
	createdAt := int64(orderObject["created"].(float64))
	dateString := time.Unix(createdAt, 0).Format("2006-01-02")
	dateStringFormatted := fmt.Sprint(strings.Split(dateString, "-")[1], "/", strings.Split(dateString, "-")[2], "/", strings.Split(dateString, "-")[0])

	customerId, ok := orderObject["customer"].(string)
	if !ok {
		logger.Error("CustomerID could not be retrieved from stripe  webhook event")
	}
	orderId, ok := orderObject["id"].(string)
	if !ok {
		logger.Error("OrderID could not be retrieved from stripe  webhook event")
	}
	order, err := fetchOrderById(orderId)
	if err != nil {
		logger.Error("Could not get order from id ", err)
	}

	stripeCustomer, _ := fetchCustomerFromId(customerId)
	card, _ := fetchCard(customerId, stripeCustomer.DefaultSource.ID)

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
	emailInfo.Shipping.Address = order.Shipping.Address.Line1
	emailInfo.Shipping.City = order.Shipping.Address.City
	emailInfo.Shipping.State = order.Shipping.Address.State
	emailInfo.Shipping.Zip = order.Shipping.Address.Zip

	// following is values from meta for easypost
	shippingId, ok := order.Meta["shipment_id"]
	if ok {
		easypostShipment, _ := fetchShipmentFromId(shippingId)
		emailInfo.Shipping.EstimatedArrival = "4-7"
		emailInfo.Shipping.TrackingCarrier = easypostShipment.Tracker.Carrier
		emailInfo.Shipping.TrackingNumber = easypostShipment.TrackingCode
		emailInfo.Shipping.TrackingUrl = easypostShipment.Tracker.PublicURL
	}

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
	firstName, _ := nameParser(order.Shipping.Name)
	emailInfo.FirstName = firstName
	emailInfo.OrderNumber = strings.Replace(orderObject["id"].(string), "or_", "", -1)
	emailInfo.OrderTotal = orderObject["amount"].(float64) / 100.0

	return emailInfo, nil

}

func MailgunSendEmail(email Email, tag string, deliveryTime time.Time) (res string, err error) {
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
	message.SetDeliveryTime(deliveryTime)
	message.AddTag(tag)
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
