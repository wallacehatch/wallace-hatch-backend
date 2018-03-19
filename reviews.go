package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	collectionName = "reviews"
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

type productReview struct {
	ProductId     string    `json:"product_id" bson:"product_id"`
	StarRating    float32   `json:"star_rating" bson:"star_rating"`
	ReviewTitle   string    `json:"review_title" bson:"review_title"`
	ReviewMessage string    `json:"review_message" bson:"review_message"`
	CustomerId    string    `json:"customer_id" bson:"customer_id"`
	CustomerName  string    `json:"customer_name" bson:"customer_name"`
	CustomerEmail string    `json:"customer_email" bson:"customer_email"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpVotes       int       `json:"up_votes" bson:"up_votes"`
	DownVotes     int       `json:"down_votes" bson:"down_votes"`
}

type productReviewRequest struct {
	ProductId     string  `json:"product_id" bson:"product_id"`
	StarRating    float32 `json:"star_rating" bson:"star_rating"`
	ReviewTitle   string  `json:"review_title" bson:"review_title"`
	ReviewMessage string  `json:"review_message" bson:"review_message"`
	CustomerEmail string  `json:"customer_email" bson:"customer_email"`
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

	customer, _ := getCustomerFromEmail(productReviewReq.CustomerEmail)

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
	pr.UpVotes = 0
	pr.DownVotes = 0

	err = db.C(collectionName).Insert(pr)
	if err != nil {
		logger.Error("Error inserting product review", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	js, _ := json.Marshal(pr)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

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
	js, _ := json.Marshal(productReviews)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

// Finds existing budget tied to user, if none exists it creates an empty budget and will return the newly created budget
func readProductReviews(productId string, db *mgo.Database) ([]productReview, error) {

	prs := make([]productReview, 0)
	err := db.C(collectionName).Find(bson.M{"product_id": productId}).All(&prs)
	if err != nil {
		logger.Error("Erorr getting all reviews ", err)
	}
	return prs, err

}

func validateReviewHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("validating revuew")
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
		return false, errors.New("No customer found with provided email")
	}
	reviews, _ := readProductReviews(productId, db)
	prevPurchased := doesCustomerContainPastOrder(customer.ID, productId)
	if !prevPurchased {
		return false, errors.New("Customer with provided email has never purchased product before")
	}
	for _, rev := range reviews {
		if rev.CustomerEmail == email {
			return false, errors.New("Customer with provided email has already left a review for this product")
		}
	}
	return true, nil
}
