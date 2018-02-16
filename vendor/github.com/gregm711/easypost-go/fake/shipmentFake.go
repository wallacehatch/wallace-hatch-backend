package fake

var ShipmentCreate string
var ShipmentBuy string

func init() {
	ShipmentCreate = `
	{
	  "id": "shp_vN9h7XLn",
	  "object": "Shipment",
	  "mode": "test",
	  "to_address": {
	    "id": "adr_zMlCRtmt",
	    "object": "Address",
	    "name": "Dr. Steve Brule",
	    "company": null,
	    "street1": "179 N Harbor Dr",
	    "street2": null,
	    "city": "Redondo Beach",
	    "state": "CA",
	    "zip": "90277",
	    "country": "US",
	    "phone": "4153334444",
	    "mode": "test",
	    "carrier_facility": null,
	    "residential": null,
	    "email": "dr_steve_brule@gmail.com",
	    "created_at": "2013-04-22T05:39:56Z",
	    "updated_at": "2013-04-22T05:39:56Z"
	  },
	  "from_address": {
	    "id": "adr_VgoLT6Ex",
	    "object": "Address",
	    "name": "EasyPost",
	    "company": null,
	    "street1": "417 Montgomery Street",
	    "street2": "5th Floor",
	    "city": "San Francisco",
	    "state": "CA",
	    "zip": "94104",
	    "country": "US",
	    "phone": "4153334444",
	    "email": "support@easypost.com",
	    "mode": "test",
	    "carrier_facility": null,
	    "residential": null,
	    "created_at": "2013-04-22T05:39:57Z",
	    "updated_at": "2013-04-22T05:39:57Z"
	  },
	  "parcel": {
	    "id": "prcl_SI55pjx8",
	    "object": "Parcel",
	    "length": 20.2,
	    "width": 10.9,
	    "height": 5.0,
	    "predefined_package": null,
	    "weight": 140.8,
	    "created_at": "2013-04-22T05:39:57Z",
	    "updated_at": "2013-04-22T05:39:57Z"
	  },
	  "customs_info": {
	    "id": "cstinfo_iNAizey5",
	    "object": "CustomsInfo",
	    "created_at": "2013-04-22T05:39:57Z",
	    "updated_at": "2013-04-22T05:39:57Z",
	    "contents_explanation": null,
	    "contents_type": "merchandise",
	    "customs_certify": false,
	    "customs_signer": null,
	    "eel_pfc": null,
	    "non_delivery_option": "return",
	    "restriction_comments": null,
	    "restriction_type": "none",
	    "customs_items": [
	      {
	        "id": "cstitem_9eDIbaDR",
	        "object": "CustomsItem",
	        "description": "Many, many EasyPost stickers.",
	        "hs_tariff_number": "123456",
	        "origin_country": "US",
	        "quantity": 1,
	        "value": 879,
	        "weight": 140,
	        "created_at": "2013-04-22T05:39:57Z",
	        "updated_at": "2013-04-22T05:39:57Z"
	      }
	    ]
	  },
	  "rates": [
	    {
	      "id": "rate_nyCb6ubX",
	      "object": "Rate",
	      "carrier_account_id": "ca_12345678",
	      "service": "FirstClassPackageInternationalService",
	      "rate": "9.50",
	      "carrier": "USPS",
	      "shipment_id": "shp_vN9h7XLn",
	      "delivery_days": 4,
	      "delivery_date": "2013-04-26T05:40:57Z",
	      "delivery_date_guaranteed": false,
	      "created_at": "2013-04-22T05:40:57Z",
	      "updated_at": "2013-04-22T05:40:57Z"
	    },
	    {
	      "id": "rate_uJh8iO2n",
	      "object": "Rate",
	      "carrier_account_id": "ca_12345678",
	      "service": "PriorityMailInternational",
	      "rate": "27.40",
	      "carrier": "USPS",
	      "shipment_id": "shp_vN9h7XLn",
	      "delivery_days": 2,
	      "delivery_date": "2013-04-24T05:40:57Z",
	      "delivery_date_guaranteed": false,
	      "created_at": "2013-04-22T05:40:57Z",
	      "updated_at": "2013-04-22T05:40:57Z"
	    },
	    {
	      "id": "rate_oZqapNpE",
	      "object": "Rate",
	      "carrier_account_id": "ca_12345678",
	      "service": "ExpressMailInternational",
	      "rate": "35.48",
	      "carrier": "USPS",
	      "shipment_id": "shp_vN9h7XLn",
	      "delivery_days": 1,
	      "delivery_date": "2013-04-23T05:40:57Z",
	      "delivery_date_guaranteed": true,
	      "created_at": "2013-04-22T05:40:57Z",
	      "updated_at": "2013-04-22T05:40:57Z"
	    }
	  ],
	  "scan_form": null,
	  "selected_rate": null,
	  "postage_label": null,
	  "tracking_code": null,
	  "refund_status": null,
	  "insurance": null,
	  "created_at": "2013-04-22T05:40:57Z",
	  "updated_at": "2013-04-22T05:40:57Z"
	}
	`

	ShipmentBuy = `
	{
	  "batch_message": null,
	  "batch_status": null,
	  "created_at": "2013-11-08T15:50:00Z",
	  "customs_info": null,
	  "from_address": {
	    "city": "San Francisco",
	    "company": null,
	    "country": "US",
	    "created_at": "2013-11-08T15:49:59Z",
	    "email": null,
	    "id": "adr_faGeob9S",
	    "mode": "test",
	    "name": "EasyPost",
	    "object": "Address",
	    "phone": "415-379-7678",
	    "state": "CA",
	    "street1": "417 Montgomery Street",
	    "street2": "5th Floor",
	    "updated_at": "2013-11-08T15:49:59Z",
	    "zip": "94104"
	  },
	  "id": "shp_vN9h7XLn",
	  "insurance": null,
	  "is_return": false,
	  "mode": "test",
	  "object": "Shipment",
	  "parcel": {
	    "created_at": "2013-11-08T15:49:59Z",
	    "height": null,
	    "id": "prcl_TXCreKJp",
	    "length": null,
	    "mode": "test",
	    "object": "Parcel",
	    "predefined_package": "UPSLetter",
	    "updated_at": "2013-11-08T15:49:59Z",
	    "weight": 3.0,
	    "width": null
	  },
	  "postage_label": {
	    "created_at": "2013-11-08T20:57:32Z",
	    "id": "pl_2Np7asmw",
	    "integrated_form": "none",
	    "label_date": "2013-11-08T20:57:32Z",
	    "label_epl2_url": null,
	    "label_file_type": "image/png",
	    "label_pdf_url": "http://assets.geteasypost.com/postage_labels/label_pdfs/Z7t7o2.pdf",
	    "label_resolution": 200,
	    "label_size": "4x7",
	    "label_type": "default",
	    "label_url": "http://assets.geteasypost.com/postage_labels/labels/lUoagDx.png",
	    "label_zpl_url": null,
	    "object": "PostageLabel",
	    "updated_at": "2013-11-08T21:11:14Z"
	  },
	  "rates": [
	    {
	      "carrier": "UPS",
	      "created_at": "2013-11-08T15:50:02Z",
	      "currency": "USD",
	      "id": "rate_vWwNuK8z",
	      "object": "Rate",
	      "rate": "30.44",
	      "service": "NextDayAir",
	      "shipment_id": "shp_w1xi76n4",
	      "updated_at": "2013-11-08T15:50:02Z"
	    }, {
	      "carrier": "UPS",
	      "created_at": "2013-11-08T15:50:02Z",
	      "currency": "USD",
	      "id": "rate_kqsS35Cx",
	      "object": "Rate",
	      "rate": "60.28",
	      "service": "NextDayAirEarlyAM",
	      "shipment_id": "shp_w1xi76n4",
	      "updated_at": "2013-11-08T15:50:02Z"
	    }
	  ],
	  "reference": null,
	  "refund_status": null,
	  "scan_form": null,
	  "selected_rate": {
	    "carrier": "UPS",
	    "created_at": "2013-11-08T15:50:02Z",
	    "currency": "USD",
	    "id": "rate_vWwNuK8z",
	    "object": "Rate",
	    "rate": "30.44",
	    "service": "NextDayAir",
	    "shipment_id": "shp_w1xi76n4",
	    "updated_at": "2013-11-08T15:50:02Z"
	  },
	  "status": "unknown",
	  "to_address": {
	    "city": "Redondo Beach",
	    "company": null,
	    "country": "US",
	    "created_at": "2013-11-08T15:49:58Z",
	    "email": "dr_steve_brule@gmail.com",
	    "id": "adr_EyDCpoii",
	    "mode": "test",
	    "name": "Dr. Steve Brule",
	    "object": "Address",
	    "phone": null,
	    "state": "CA",
	    "street1": "179 N Harbor Dr",
	    "street2": null,
	    "updated_at": "2013-11-08T15:49:58Z",
	    "zip": "90277"
	  },
	  "tracker": {
	    "created_at": "2013-11-08T20:57:32Z",
	    "id": "trk_k7oNSa1Q",
	    "mode": "test",
	    "object": "Tracker",
	    "shipment_id": "shp_w1xi76n4",
	    "status": "unknown",
	    "tracking_code": "1ZE6A4850190733810",
	    "tracking_details": [ ],
	    "updated_at": "2013-11-08T20:58:26Z"
	    "public_url": "https://track.easypost.com/djE7...",
	  },
	  "tracking_code": "1ZE6A4850190733810",
	  "updated_at": "2013-11-08T20:58:26Z"
	}
	`
}
