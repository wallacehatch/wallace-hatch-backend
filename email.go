package main

import (
	"bytes"
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
	"html/template"
	"os"
	"strings"
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

func constructEmail() {
	email := make(map[string]interface{})
	bufferBytes := bytes.Buffer{}
	tmpl, err := template.ParseFiles("email-templates/account-ready.html")
	if err != nil {
		fmt.Println("error opening template ", err)
		return
	}
	info := EmailInformation{}
	if err := tmpl.Execute(&bufferBytes, info); err != nil {
		fmt.Println("error executing html ", err)
		return
	}

	email["subject"] = "Your data is ready!"
	email["plain_text"] = "data"
	email["to"] = "greg711miller@gmail.com"
	email["from"] = "info@wallacehatch.com"
	email["html"] = bufferBytes.String()
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
	trimmed := strings.Trim(id, "<")
	result := strings.Split(trimmed, "@")[0:1]
	joined := strings.Join(result, "")
	finalId := strings.Join(strings.Split(joined, "."), "")
	if err != nil {
		fmt.Println("MAILGUN ERROR ", err)
		return "", err
	}
	return finalId, err
}
