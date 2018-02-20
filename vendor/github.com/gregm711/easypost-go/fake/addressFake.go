package fake

var EasypostFakeAddress string

func init() {
	EasypostFakeAddress = `
		{
		  "id": "adr_3dc0df9695654f1f8a788e1c872a0769",
		  "object": "Address",
		  "created_at": "2015-12-22T19:28:06Z",
		  "updated_at": "2015-12-22T19:28:06Z",
		  "name": null,
		  "company": "EasyPost",
		  "street1": "417 MONTGOMERY ST",
		  "street2": "FL 5",
		  "city": "SAN FRANCISCO",
		  "state": "CA",
		  "zip": "94104-110",
		  "country": "US",
		  "phone": "4151234567",
		  "email": null,
		  "mode": "test",
		  "carrier_facility": null,
		  "residential": null,
		  "federal_tax_id": null,
		  "state_tax_id": null,
		  "verifications": {
		    "delivery": {
		      "success": true,
		      "errors": []
		    }
		  }
		}
	`
}
