package main

import (
	"encoding/json"
	"fmt"
	"github.com/sfreiberg/gotwilio"
	"os"
	"unicode"
)

var twillioSID string
var twillioAuthToken string
var twillioPhoneNumber string

func init() {
	twillioSID = os.Getenv("TWILIO_ACCOUNT_SID")
	twillioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	twillioPhoneNumber = os.Getenv("TWILIO_PHONE_NUMBER")

}

func sendMessage(toPhoneNumber string, message string) error {
	toPhoneNumber = "4403966613" // only send to greg's cell until we know shit is working
	toPhoneNumber = cleanPhoneNumber(toPhoneNumber)
	twilio := gotwilio.NewTwilioClient(twillioSID, twillioAuthToken)
	response, _, err := twilio.SendSMS(twillioPhoneNumber, toPhoneNumber, message, "", "")
	if err != nil {
		logger.Error("Error sending message through twilio ", err)
	}
	_, err = json.Marshal(response)
	if err != nil {
		logger.Error("Error marashling twilio response ", err)
	}
	return err
}

// ex : go from  +1 (440) 396-6613  -> +14403966613
func cleanPhoneNumber(s string) string {
	cleanedPhoneNumber := "+"
	for _, char := range s {
		if unicode.IsNumber(char) {
			cleanedPhoneNumber = fmt.Sprint(cleanedPhoneNumber, string(char))
		}
	}
	return cleanedPhoneNumber
}
