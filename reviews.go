package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	collectionName = "reviews"
	starFilledUrl  = "https://s3.us-east-2.amazonaws.com/wallace-hatch/fa-star-filled-14.png"
	starWhiteUrl   = "https://s3.us-east-2.amazonaws.com/wallace-hatch/fa-star-white.png"
)

func db() *mgo.Database {
	url := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_NAME")
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	dialInfo, err := mgo.ParseURL(url)

	if err != nil {
		logger.Error("Can't create dial info for db:", err)
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Error("Can't create db session:", err)
	}
	defer session.Close()

	return session.Clone().DB(dbName)
}

func createProductReviewHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var productReviewReq productReviewRequest
	err := decoder.Decode(&productReviewReq)
	if err != nil {
		logger.Error("Error decoding product review request", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	_, err = validateCustomerReview(productReviewReq.CustomerEmail, productReviewReq.ProductId)
	if err != nil {
		logger.Error("Customer is not validated: ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return

	}

	customer, err := getCustomerFromEmail(productReviewReq.CustomerEmail)
	if err != nil {
		logger.Error("Error getting customer ", err)
	}

	db := db()
	defer db.Session.Close()
	pr := productReview{}
	pr.ProductId = productReviewReq.ProductId
	pr.StarRating = productReviewReq.StarRating
	pr.ReviewTitle = productReviewReq.ReviewTitle
	pr.ReviewMessage = productReviewReq.ReviewMessage
	pr.CustomerId = customer.ID
	pr.CustomerName = customer.Meta["name"]
	pr.CustomerEmail = customer.Email
	pr.CreatedAt = time.Now()
	pr.FriendRecommendation = productReviewReq.FriendRecommendation
	pr.BrandRating = productReviewReq.BrandRating
	pr.BrandRatingMessage = productReviewReq.BrandRatingMessage

	err = db.C(collectionName).Insert(pr)
	if err != nil {
		logger.Error("Error inserting product review", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}

	sendReviewEmail(pr)
	js, _ := json.Marshal(pr)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

func sendReviewEmail(rev productReview) {

	bufferBytes := bytes.Buffer{}

	emailItem := EmailItemInformation{}
	emailItems := make([]EmailItemInformation, 0)
	product, _ := fetchProductById(rev.ProductId)

	emailItem.ImageUrl = product.Images[0]
	emailItem.Name = product.Name
	emailInfo := EmailInformation{}
	emailInfo.FirstName, _ = nameParser(rev.CustomerName)
	emailItems = append(emailItems, emailItem)
	emailInfo.Items = emailItems
	emailInfo.ReviewTitle = rev.ReviewTitle
	emailInfo.ReviewMessage = rev.ReviewMessage

	switch rev.StarRating {
	case 1.0:
		emailInfo.StarOneUrl = starFilledUrl
		emailInfo.StarTwoUrl = starWhiteUrl
		emailInfo.StarThreeUrl = starWhiteUrl
		emailInfo.StarFourUrl = starWhiteUrl
		emailInfo.StarFiveUrl = starWhiteUrl
	case 2.0:
		emailInfo.StarOneUrl = starFilledUrl
		emailInfo.StarTwoUrl = starFilledUrl
		emailInfo.StarThreeUrl = starWhiteUrl
		emailInfo.StarFourUrl = starWhiteUrl
		emailInfo.StarFiveUrl = starWhiteUrl
	case 3.0:
		emailInfo.StarOneUrl = starFilledUrl
		emailInfo.StarTwoUrl = starFilledUrl
		emailInfo.StarThreeUrl = starFilledUrl
		emailInfo.StarFourUrl = starWhiteUrl
		emailInfo.StarFiveUrl = starWhiteUrl
	case 4.0:
		emailInfo.StarOneUrl = starFilledUrl
		emailInfo.StarTwoUrl = starFilledUrl
		emailInfo.StarThreeUrl = starFilledUrl
		emailInfo.StarFourUrl = starFilledUrl
		emailInfo.StarFiveUrl = starWhiteUrl
	case 5.0:
		emailInfo.StarOneUrl = starFilledUrl
		emailInfo.StarTwoUrl = starFilledUrl
		emailInfo.StarThreeUrl = starFilledUrl
		emailInfo.StarFourUrl = starFilledUrl
		emailInfo.StarFiveUrl = starFilledUrl
	}
	tmpl, err := template.ParseFiles("email-templates/review.html")
	if err != nil {
		logger.Error("error opening template ", err)
	}
	if err := tmpl.Execute(&bufferBytes, emailInfo); err != nil {
		logger.Error("error executing html ", err)
	}
	email := Email{}
	email.Subject = "Your Review is live!"
	email.PlainText = "Your review is live!"
	email.From = emailSender
	email.To = rev.CustomerEmail
	email.Html = bufferBytes.String()
	MailgunSendEmail(email, reviewTag, time.Now())

}

func fetchProductReviewsHandler(w http.ResponseWriter, r *http.Request) {
	db := db()
	defer db.Session.Close()
	vars := mux.Vars(r)
	productReviews, err := readProductReviews(vars["key"], db)
	if err != nil {
		logger.Error("Error retireving product reviews", err)
		respondErrorJson(err, http.StatusBadRequest, w)
	}

	productReviewsResponse := make([]productReviewResp, 0)
	js, _ := json.Marshal(productReviews)
	json.Unmarshal(js, &productReviewsResponse)

	for index, review := range productReviewsResponse {
		customerReviews, _ := readCustomerReviews(review.CustomerId, db)
		productReviewsResponse[index].CustomerReviews = len(customerReviews)
	}

	jsResp, _ := json.Marshal(productReviewsResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResp)
	return
}

func readCustomerReviews(customerId string, db *mgo.Database) ([]productReview, error) {
	prs := make([]productReview, 0)
	err := db.C(collectionName).Find(bson.M{"customer_id": customerId}).All(&prs)
	if err != nil {
		logger.Error("Erorr getting reviews by customer id ", err)
	}
	return prs, err

}

// Finds existing budget tied to user, if none exists it creates an empty budget and will return the newly created budget
func readProductReviews(productId string, db *mgo.Database) ([]productReview, error) {

	prs := make([]productReview, 0)
	err := db.C(collectionName).Find(bson.M{"product_id": productId}).All(&prs)
	if err != nil {
		logger.Error("Erorr getting reviews by product id reviews ", err)
	}
	return prs, err

}

func validateReviewHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var valReviewReq validateReviewRequest
	err := decoder.Decode(&valReviewReq)
	if err != nil {
		logger.Error("Error decoding validate review request", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	logger.Info("valReviewReq.CustomerEmail")
	validatedReview, err := validateCustomerReview(valReviewReq.CustomerEmail, valReviewReq.ProductId)
	if err != nil {
		logger.Error("Customer is probably not validated: ", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return

	}
	result := make(map[string]interface{})
	result["verified_buyer"] = validatedReview
	js, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	logger.Info("was verified!")
	return

}

// determines if user is qualifed to create and leave review
func validateCustomerReview(email string, productId string) (bool, error) {
	db := db()
	defer db.Session.Close()
	customer, err := getCustomerFromEmail(email)
	if err != nil {
		logger.Error("No customers found for this email.", email)
		return false, errors.New("No customers found for this email.")
	}
	reviews, _ := readProductReviews(productId, db)
	prevPurchased := doesCustomerContainPastOrder(customer.ID, productId)
	if !prevPurchased {
		logger.Error("This email has never purchased product before ", email)
		return false, errors.New("Incorrect order email, this email has never purchased this product before.")
	}
	for _, rev := range reviews {
		if rev.CustomerEmail == email {
			logger.Error("This email has already left a review for this product.", email)
			return false, errors.New("This email has already left a review for this product.")
		}
	}
	return true, nil
}
