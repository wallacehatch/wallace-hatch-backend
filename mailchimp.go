package main

import (
	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
)

const (
	mailchimpNewsletterListId = "ae8e38bbe3"
)

func addToMailchimpNewsletter(email string, company string, message string) (members.Member, error) {
	err := mailchimp.SetKey(mailchimpAPIKey)
	if err != nil {
		logger.Error("error connecting to mailchimp api", err)
		return members.Member{}, err
	}
	status := "subscribed"
	mergeFields := make(map[string]interface{})
	mergeFields["COMPANY"] = company
	mergeFields["MESSAGE"] = message
	params := &members.NewParams{}
	params.EmailAddress = email
	params.MergeFields = mergeFields
	params.Status = members.Status(status)
	member, err := members.New(mailchimpNewsletterListId, params)
	if err != nil {
		logger.Error("error adding to mailchimp", err)
		return members.Member{}, err
	}
	return *member, err
}
