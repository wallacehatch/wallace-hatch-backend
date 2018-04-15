package main

import (
	"time"
)

type accountRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	AcceptTerms bool   `json:"acceptTerms"`
}

type cardRequest struct {
	Number string `json:"cardNumber"`
	Exp    string `json:"exp"`
	CVC    string `json:"cvc"`
}

type orderRequest struct {
	Items []itemRequest `json:"items"`
}

type couponSubmitRequest struct {
	Email string `json:"email"`
}

type itemRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type shippingRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	StreetNumber string `json:"streetNumber"`
	StreetName   string `json:"streetName"`
	AptSuite     string `json:"aptSuite"`
	Company      string `json:"company"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Zip          string `json:"zip"`
}

type Response struct {
	Text   string `json:"text"`
	Status int    `json:"status"`
}

type completeOrderRequest struct {
	Account     accountRequest         `json:"account"`
	Card        cardRequest            `json:"cardInfo"`
	Order       orderRequest           `json:"cart"`
	Shipping    shippingRequest        `json:"shipping"`
	Coupon      string                 `json:"coupon"`
	GooglePlace map[string]interface{} `json:"google_place"`
}

type Email struct {
	From      string `json:"from" bson:"from"`
	To        string `json:"to" bson:"to"`
	Subject   string `json:"subject" bson:"subject"`
	Html      string `json:"html" bson:"html"`
	PlainText string `json:"plain_text" bson:"plain_text"`
}

type EmailItemInformation struct {
	Name     string
	Size     string
	Color    string
	Style    string
	ImageUrl string
	Quantity int
	Price    float64
}

type EmailShippingInformation struct {
	Address          string
	City             string
	State            string
	Zip              string
	TrackingNumber   string
	TrackingCarrier  string
	EstimatedArrival string // dont need to get fancy ,just set to 4-7 days for now
	TrackingUrl      string
}

type EmailInformation struct {
	To             string
	From           string
	FirstName      string
	OrderNumber    string
	OrderDate      string
	OrderTotal     float64
	Items          []EmailItemInformation
	CardType       string
	CardImageUrl   string
	CardMask       string
	Shipping       EmailShippingInformation
	NumItems       int
	NumItemsMinus  int
	CouponCode     string
	CouponDiscount int
	StarOneUrl     string
	StarTwoUrl     string
	StarThreeUrl   string
	StarFourUrl    string
	StarFiveUrl    string
	ReviewMessage  string
	ReviewTitle    string
}

type ResponseError struct {
	ErrorMsg string `json:"error_message"`
	Status   int    `json:"status"`
}

type idsReqeust struct {
	Ids []string `json:"product_ids"`
}

type validateReviewRequest struct {
	ProductId     string `json:"product_id"`
	CustomerEmail string `json:"customer_email"`
}

type instagramCommentResp struct {
	Data []struct {
		ID   string `json:"id"`
		From struct {
			ID             string `json:"id"`
			Username       string `json:"username"`
			FullName       string `json:"full_name"`
			ProfilePicture string `json:"profile_picture"`
		} `json:"from"`
		Text        string `json:"text"`
		CreatedTime string `json:"created_time"`
	} `json:"data"`
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
}

type instagramMediaResp struct {
	Data struct {
		Type         string `json:"type"`
		UsersInPhoto []struct {
			User struct {
				Username       string `json:"username"`
				FullName       string `json:"full_name"`
				ID             string `json:"id"`
				ProfilePicture string `json:"profile_picture"`
			} `json:"user"`
			Position struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"position"`
		} `json:"users_in_photo"`
		Filter   string        `json:"filter"`
		Tags     []interface{} `json:"tags"`
		Comments struct {
			Count int `json:"count"`
		} `json:"comments"`
		Caption interface{} `json:"caption"`
		Likes   struct {
			Count int `json:"count"`
		} `json:"likes"`
		Link string `json:"link"`
		User struct {
			Username       string `json:"username"`
			FullName       string `json:"full_name"`
			ProfilePicture string `json:"profile_picture"`
			ID             string `json:"id"`
		} `json:"user"`
		CreatedTime string `json:"created_time"`
		Images      struct {
			LowResolution struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"low_resolution"`
			Thumbnail struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnail"`
			StandardResolution struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"standard_resolution"`
		} `json:"images"`
		ID       string      `json:"id"`
		Location interface{} `json:"location"`
	} `json:"data"`
}

type easypostWebhook struct {
	Result struct {
		ID              string      `json:"id"`
		Object          string      `json:"object"`
		Mode            string      `json:"mode"`
		TrackingCode    string      `json:"tracking_code"`
		Status          string      `json:"status"`
		StatusDetail    string      `json:"status_detail"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		SignedBy        interface{} `json:"signed_by"`
		Weight          interface{} `json:"weight"`
		EstDeliveryDate time.Time   `json:"est_delivery_date"`
		ShipmentID      string      `json:"shipment_id"`
		Carrier         string      `json:"carrier"`
		TrackingDetails []struct {
			Object           string      `json:"object"`
			Message          string      `json:"message"`
			Status           string      `json:"status"`
			StatusDetail     string      `json:"status_detail"`
			Datetime         time.Time   `json:"datetime"`
			Source           string      `json:"source"`
			CarrierCode      interface{} `json:"carrier_code"`
			TrackingLocation struct {
				Object  string      `json:"object"`
				City    interface{} `json:"city"`
				State   interface{} `json:"state"`
				Country interface{} `json:"country"`
				Zip     interface{} `json:"zip"`
			} `json:"tracking_location"`
		} `json:"tracking_details"`
		CarrierDetail struct {
			Object                 string      `json:"object"`
			Service                string      `json:"service"`
			ContainerType          interface{} `json:"container_type"`
			EstDeliveryDateLocal   interface{} `json:"est_delivery_date_local"`
			EstDeliveryTimeLocal   interface{} `json:"est_delivery_time_local"`
			OriginLocation         string      `json:"origin_location"`
			OriginTrackingLocation struct {
				Object  string      `json:"object"`
				City    string      `json:"city"`
				State   string      `json:"state"`
				Country interface{} `json:"country"`
				Zip     string      `json:"zip"`
			} `json:"origin_tracking_location"`
			DestinationLocation         string      `json:"destination_location"`
			DestinationTrackingLocation interface{} `json:"destination_tracking_location"`
			GuaranteedDeliveryDate      interface{} `json:"guaranteed_delivery_date"`
			AlternateIdentifier         interface{} `json:"alternate_identifier"`
			InitialDeliveryAttempt      interface{} `json:"initial_delivery_attempt"`
		} `json:"carrier_detail"`
		Finalized bool          `json:"finalized"`
		IsReturn  bool          `json:"is_return"`
		PublicURL string        `json:"public_url"`
		Fees      []interface{} `json:"fees"`
	} `json:"result"`
	Description        string `json:"description"`
	Mode               string `json:"mode"`
	PreviousAttributes struct {
		Status string `json:"status"`
	} `json:"previous_attributes"`
	CreatedAt     time.Time     `json:"created_at"`
	PendingUrls   []string      `json:"pending_urls"`
	CompletedUrls []interface{} `json:"completed_urls"`
	UpdatedAt     time.Time     `json:"updated_at"`
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	Status        string        `json:"status"`
	Object        string        `json:"object"`
}

type productReviewResp struct {
	ProductId            string    `json:"product_id" bson:"product_id"`
	StarRating           float32   `json:"star_rating" bson:"star_rating"`
	ReviewTitle          string    `json:"review_title" bson:"review_title"`
	ReviewMessage        string    `json:"review_message" bson:"review_message"`
	CustomerReviews      int       `json:"customer_reviews" bson:"customer_reviews"`
	CustomerId           string    `json:"customer_id" bson:"customer_id"`
	CustomerName         string    `json:"customer_name" bson:"customer_name"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at"`
	FriendRecommendation bool      `json:"friend_recommendation" bson:"friend_recommendation"`
}

type productReview struct {
	ProductId            string    `json:"product_id" bson:"product_id"`
	StarRating           float32   `json:"star_rating" bson:"star_rating"`
	ReviewTitle          string    `json:"review_title" bson:"review_title"`
	ReviewMessage        string    `json:"review_message" bson:"review_message"`
	CustomerId           string    `json:"customer_id" bson:"customer_id"`
	CustomerName         string    `json:"customer_name" bson:"customer_name"`
	CustomerEmail        string    `json:"customer_email" bson:"customer_email"`
	CreatedAt            time.Time `json:"created_at" bson:"created_at"`
	FriendRecommendation bool      `json:"friend_recommendation" bson:"friend_recommendation"`
	BrandRating          float32   `json:"brand_rating" bson:"brand_rating"`
	BrandRatingMessage   string    `json:"brand_rating_message" bson:"brand_rating_message"`
}

type productReviewRequest struct {
	ProductId            string  `json:"product_id" bson:"product_id"`
	StarRating           float32 `json:"star_rating" bson:"star_rating"`
	ReviewTitle          string  `json:"review_title" bson:"review_title"`
	ReviewMessage        string  `json:"review_message" bson:"review_message"`
	CustomerEmail        string  `json:"customer_email" bson:"customer_email"`
	FriendRecommendation bool    `json:"friend_recommendation" bson:"friend_recommendation"`
	BrandRating          float32 `json:"brand_rating" bson:"brand_rating"`
	BrandRatingMessage   string  `json:"brand_rating_message" bson:"brand_rating_message"`
}
