package main

import (
	"github.com/gregm711/easypost-go/easypost"
	"github.com/stripe/stripe-go"

	"os"
	"strconv"
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
	itemWieght      = 1.0 // lbs
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
	shipment, _ := createShipment(toAdd, fromAdd, parcel)
	lowestPriceShipment, _ := selectLowestShipmentRate(shipment)
	boughtShipment, err := buyShipment(lowestPriceShipment)
	if err != nil {
		return
	}
	logger.Info(boughtShipment.Tracker)
	logger.Info(boughtShipment.Tracker.EstDeliveryDate.)

	updatedMeta := map[string]string{"postage_label": boughtShipment.PostageLabel.LabelURL, "tracking_code": boughtShipment.TrackingCode}
	updateOrderMeta(order.ID, updatedMeta)
	logger.Info("updated order")

}

func createFromAdd() (easypost.Address, error) {
	address := easypost.Address{
		Name:    "Wallce Hatch",
		Street1: "364 East Northwood Avenue",
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

	totalweight := 0.0
	for _, item := range order.Items {
		totalweight = totalweight + float64(itemWieght*item.Quantity)
	}

	parcel := easypost.Parcel{
		PredefinedPackage: uspsPackageType,
		Weight:            float32(totalweight),
	}
	err := parcel.Create()
	if err != nil {
		logger.Error("Error creating parcel", err)
	}
	return parcel, err

}

func createShipment(toAddress easypost.Address, fromAddress easypost.Address, parcel easypost.Parcel) (easypost.Shipment, error) {
	shipment := easypost.Shipment{
		ToAddress:     toAddress,
		FromAddress:   fromAddress,
		ReturnAddress: fromAddress,
		Parcel:        parcel,
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