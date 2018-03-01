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
	CustomerId    string  `json:"customer_id" bson:"customer_id"`
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
	customer, _ := fetchCustomerFromId(productReviewReq.CustomerId)

	prevPurchased := doesCustomerContainPastOrder(customer.ID, productReviewReq.ProductId)

	if !prevPurchased {
		err := errors.New("User has not purchased this product before, and therefore cannnot leave a review")
		logger.Error("User has not purchased this product before", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
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
