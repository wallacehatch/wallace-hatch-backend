package main

import (
	"fmt"
	"github.com/gregm711/easypost-go/easypost"
	"github.com/stripe/stripe-go"
	"os"
	"strconv"
	"time"
)

/*
Backend steps for ordering
1: Customer is created
2: Card is created and attached to customer
3: order is created and attached to customer
4: Order is charged by defualt card
5: Easy post from and to  Addresses is created
6: Easy post parcel is created based off of order
7: Easy post shipment rate ID is created
8: Easy post shipment is purchased, returning a tracking code url for paid label image, and url for branded tracking page
9: Label and tracking ID is attached to order to be used when fulfilment occurs
10: Webhooks are processed from easy post and SMS messages are sent out accordingly
*/

const (
	itemWieght      = 9.0 // ounces
	uspsPackageType = "parcel"
)

func init() {
	easypost.SetApiKey(os.Getenv("EASY_POST_KEY"))
	easypost.Request = easypost.RequestController{}
}

func easypostController(order stripe.Order) {

	toAdd, _ := createAddress(order.Shipping)
	fromAdd, _ := createFromAdd()

	parcel, _ := createParcel(order)
	shipment, _ := createShipment(toAdd, fromAdd, parcel, order.ID)
	lowestPriceShipment, _ := selectLowestShipmentRate(shipment)
	boughtShipment, err := buyShipment(lowestPriceShipment)
	if err != nil {
		return
	}
	updatedMeta := map[string]string{"postage_label_pdf_url": boughtShipment.PostageLabel.LabelPdfURL, "postage_label_url": boughtShipment.PostageLabel.LabelURL, "tracking_code": boughtShipment.TrackingCode, "shipment_id": boughtShipment.ID, "tracking_url": boughtShipment.Tracker.PublicURL}
	updateOrderMeta(order.ID, updatedMeta)
}

func createFromAdd() (easypost.Address, error) {
	address := easypost.Address{
		Name:    "Wallce Hatch",
		Street1: "629 North High Street",
		Street2: "4th Floor",
		City:    "Columbus",
		State:   "OH",
		Zip:     "43201",
		Country: "US",
		Email:   "customerservice@wallacehatch.com",
		Phone:   "4159150936",
	}
	err := address.Create()

	return address, err

}

func createAddress(shippingInfo stripe.Shipping) (easypost.Address, error) {

	address := easypost.Address{
		Name:    shippingInfo.Name,
		Street1: shippingInfo.Address.Line1,
		Street2: shippingInfo.Address.Line2,
		City:    shippingInfo.Address.City,
		State:   shippingInfo.Address.State,
		Zip:     shippingInfo.Address.Zip,
		Country: shippingInfo.Address.Country,
	}

	err := address.Create()
	if err != nil {
		logger.Error("Error creating address", err)
	}

	return address, err

}

func createParcel(order stripe.Order) (easypost.Parcel, error) {

	totalWeight := 0.0
	for _, item := range order.Items {
		totalWeight = totalWeight + float64(itemWieght*item.Quantity)
	}
	logger.Info(totalWeight)
	parcel := easypost.Parcel{
		PredefinedPackage: uspsPackageType,
		Weight:            float32(totalWeight),
	}
	err := parcel.Create()
	if err != nil {
		logger.Error("Error creating parcel", err)
	}
	return parcel, err

}

func createShipment(toAddress easypost.Address, fromAddress easypost.Address, parcel easypost.Parcel, stripeOrderId string) (easypost.Shipment, error) {
	shipment := easypost.Shipment{
		ToAddress:     toAddress,
		FromAddress:   fromAddress,
		ReturnAddress: fromAddress,
		Parcel:        parcel,
		Reference:     stripeOrderId,
	}
	err := shipment.Create()
	if err != nil {
		logger.Error("Error creating shipment", err)
	}
	return shipment, err
}

func selectLowestShipmentRate(shipment easypost.Shipment) (easypost.Shipment, error) {
	lowestRate := shipment.Rates[0]

	for index, rate := range shipment.Rates {
		price, _ := strconv.ParseFloat(rate.Rate, 64)
		currentPrice, _ := strconv.ParseFloat(lowestRate.Rate, 64)
		if price < currentPrice {
			lowestRate = shipment.Rates[index]
		}
	}
	shipment.SelectedRate = lowestRate
	return shipment, nil

}

func buyShipment(shipment easypost.Shipment) (easypost.Shipment, error) {

	err := shipment.Buy()
	if err != nil {
		logger.Error("Error buying shipment", err)
	}
	return shipment, err

}

func constructMessage(hook easypostWebhook) string {

	shortenedTrackingLink, _ := shortenUrl(hook.Result.PublicURL)
	trackingUpdates := hook.Result.TrackingDetails
	newestIndex := 0
	for index, update := range trackingUpdates {
		if update.Datetime.After(trackingUpdates[newestIndex].Datetime) {
			newestIndex = index
		}
	}
	if len(trackingUpdates) == 0 {
		return ""
	}
	mostRecentTrackingMessage := trackingUpdates[newestIndex].Message
	currentLocation := fmt.Sprint(trackingUpdates[newestIndex].TrackingLocation.City, " ", trackingUpdates[newestIndex].TrackingLocation.State)
	estimatedArrival := formatDate(hook.Result.EstDeliveryDate)
	logger.Info("for easypost webhook, most recent tracking message is : ", mostRecentTrackingMessage)

	// we know this is a juicy tracking event that the customer needs to know about
	switch mostRecentTrackingMessage {
	case "Arrived at USPS Origin Facility":
		return fmt.Sprint(mostRecentTrackingMessage, ": Your Wallace Hatch ⌚️📦 is on it's way!\n\nCurrent location 📍 ", currentLocation, "\n\nEstimated delivery 📅 ", estimatedArrival, ".\n\nTrack at ", shortenedTrackingLink)
	case "Arrived at Post Office":
		return fmt.Sprint(mostRecentTrackingMessage, ": Your Wallace Hatch ⌚️📦 is on it's way!\n\nCurrent location 📍 ", currentLocation, "\n\nEstimated delivery 📅 ", estimatedArrival, ".\n\nTrack at ", shortenedTrackingLink)
	case "Out for Delivery":
		return fmt.Sprint(mostRecentTrackingMessage, ": Your Wallace Hatch ⌚️📦 is on it's way!\n\nCurrent location 📍 ", currentLocation, "\n\nEstimated delivery 📅 ", estimatedArrival, ".\n\nTrack at ", shortenedTrackingLink)
	case "Delivered, Front Door/Porch":
		return fmt.Sprint(mostRecentTrackingMessage, ": Your Wallace Hatch ⌚️📦 has been delivered!🎉")
	}
	return ""

}

func fetchShipmentFromId(shimpentId string) (easypost.Shipment, error) {
	shipment := easypost.Shipment{}
	shipment.ID = shimpentId
	err := shipment.Get()
	if err != nil {
		logger.Error("Error fetching shipment from Id", err)
	}
	return shipment, err
}

func formatDate(date time.Time) string {
	day := ""
	switch date.Weekday().String() {
	case "Monday":
		day = "Mon"
	case "Tuesday":
		day = "Tue"
	case "Wednesday":
		day = "Wed"
	case "Thursday":
		day = "Thurs"
	case "Friday":
		day = "Tue"
	case "Saturday":
		day = "Sat"
	case "Sunday":
		day = "Sun"

	}

	return fmt.Sprint(day, ", ", date.Month().String(), " ", date.Day())

}
