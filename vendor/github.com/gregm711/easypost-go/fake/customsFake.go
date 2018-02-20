package fake

var EasypostFakeCustoms string

func init() {
	EasypostFakeCustoms = `
	{
	  "id": "cstinfo_zN5tbjEd",
	  "object": "CustomsInfo",
	  "contents_explanation": null,
	  "contents_type": "merchandise",
	  "customs_certify": true,
	  "customs_signer": "Steve Brule",
	  "eel_pfc": "NOEEI 30.37(a)",
	  "non_delivery_option": "return",
	  "restriction_comments": null,
	  "restriction_type": "none",
	  "customs_items": [{
	      "id": "cstitem_OpPpXeny",
	      "object": "CustomsItem",
	      "description": "T-Shirt",
	      "hs_tariff_number": "123456",
	      "origin_country": "US",
	      "quantity": 1,
	      "value": "10",
	      "weight": 5,
	      "created_at": "2013-04-22T07:17:51Z",
	      "updated_at": "2013-04-22T07:17:51Z"
	    }, {
	      "id": "cstitem_VklGOHPs",
	      "object": "CustomsItem",
	      "description": "Sweet shirts",
	      "hs_tariff_number": "654321",
	      "origin_country": "US",
	      "quantity": 2,
	      "value": "23",
	      "weight": 11,
	      "created_at": "2013-04-22T07:17:51Z",
	      "updated_at": "2013-04-22T07:17:51Z"
	    }
	  ],
	  "created_at": "2013-04-22T07:17:51Z",
	  "updated_at": "2013-04-22T07:17:51Z"
	}
	`
}
