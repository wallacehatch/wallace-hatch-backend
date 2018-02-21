package main

import (
	"encoding/json"
	"fmt"
	"github.com/sfreiberg/gotwilio"
	"github.com/ttacon/libphonenumber"
	"os"
	"strconv"

	"errors"
)

var twillioSID string
var twillioAuthToken string
var twillioPhoneNumber string

func init() {
	twillioSID = os.Getenv("TWILIO_ACCOUNT_SID")
	twillioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	twillioPhoneNumber = os.Getenv("TWILIO_PHONE_NUMBER")
}

func sendSMSMessage(toPhoneNumber string, message string) (*gotwilio.SmsResponse, error) {
	response := &gotwilio.SmsResponse{}
	// toPhoneNumber = "1(440) 396-6613" // only send to greg's cell until we know shit is working
	toPhoneNumber, err := cleanPhoneNumber(toPhoneNumber)
	if err != nil {
		logger.Error("Error formatting to  number", err)
		return response, err
	}
	twilio := gotwilio.NewTwilioClient(twillioSID, twillioAuthToken)
	response, exception, err := twilio.SendSMS(twillioPhoneNumber, toPhoneNumber, message, "", "")
	if exception != nil {
		logger.Error("Exception sending message through twilio", exception)
		return response, errors.New(exception.Message)
	}
	if err != nil {
		logger.Error("Error sending message through twilio ", err)
		return response, err
	}
	_, err = json.Marshal(response)
	if err != nil {
		logger.Error("Error marashling twilio response ", err)
	}
	return response, err
}

// ex : go from  +1 (440) 396-6613  -> +14403966613
func cleanPhoneNumber(s string) (string, error) {
	num, err := libphonenumber.Parse(s, "US")
	final := ""
	if err == nil {
		final = fmt.Sprint("+", strconv.Itoa(int(num.GetCountryCode())), strconv.Itoa(int(num.GetNationalNumber())))
	}
	return final, err
}
