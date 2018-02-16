# EasyPost Go Client Library
This is a library written in Go to use with your own EasyPost credentials.
More info on https://www.easypost.com

_**Note:** As far as possible, data structure in this package is reflecting the official EasyPost API structure: https://www.easypost.com/docs/api.html_

## Installation

First you need to get the package
```
go get github.com/gregm711/easypost-go
```

Then on your go project
```
import "github.com/gregm711/easypost-go/easypost"
```

## Setup

You will need to set your API key, either the test or production one, before performing any action.
Retrieve the API key associated to your account here: https://www.easypost.com/account#/api-keys

After importing the package in your project, set the API key you want to use
```
easypost.SetApiKey(__YOUR_API_KEY__)
```

Then init the request object that easypost-go will use for perform the API calls
```
easypost.Request = easypost.RequestController{}
```

_**Note:** This step is required to perform test without actually requesting the EasyPost APIs_

## Usage

Once the previous steps have been completed, you can start using the package as you would normally do in Go.

**Example**
```
var order = easypost.Order{
	ToAddress: easypost.Address {
		Street1: "417 MONTGOMERY ST",
		Street2: "FLOOR 5",
		City: "SAN FRANCISCO",
		State: "CA",
		Zip: "94104",
		Country: "US",
	},
	FromAddress: easypost.Address {
		ID: "adr_a6fd5dd822c94bdfa1e3f2d28a4dbf9b",
	}
}

order.Shipments = []easypost.Shipment{
	...
}

order.Create()
order.Buy()
```

## Errors
When performing an action with easypost objects, you can receive some errors:

```
err := order.Create()
if err != nil {
	//There was an error creating the order. Usually something not related to EasyPost
}

if order.Error != nil {
	//EasyPost returned an error when requesting order creation. Check the messages included in the Error object
}
```

Beside the generic error returned by the function, as shown in the first example, the `Error` field of the object is populated with the information EasyPost is returning.

```
//Error is an EasyPost object representing an error
type Error struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

//FieldError is an EasyPost object that defines an error in the verification
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
```
