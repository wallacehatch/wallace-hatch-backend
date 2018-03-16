package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"

	"io/ioutil"
	"net/http"
	"os"
)

var instagramAccessToken string

var watchTags = map[string]bool{"kalliorose": true, "sohohatch": true, "palermosoul": true}

const instagramApiBaseURL = "https://api.instagram.com/v1/"

func init() {
	instagramAccessToken = os.Getenv("INSTAGRAM_ACCESS_TOKEN")
}

func fetchInstagramPostInformationHandler(w http.ResponseWriter, r *http.Request) {

	var data instagramMediaData
	vars := mux.Vars(r)
	shortenedUrl := vars["key"]
	instagramMedia, err := getInstagramMediaInfo(shortenedUrl)
	if err != nil {
		logger.Error("Error with instagram media", err)
		respondErrorJson(err, http.StatusBadRequest, w)
		return
	}
	watchesInshot := make([]string, 0)
	for _, tag := range instagramMedia.Data.Tags {

		if watchTags[tag.(string)] == true {
			watchesInshot = append(watchesInshot, tag.(string))
		}
	}
	data.PictureUrl = instagramMedia.Data.Images.StandardResolution.URL
	data.WallaceHatchProfilePictureUrl = "https://instagram.fcmh1-1.fna.fbcdn.net/vp/0e5a33e58f4e512e19b038747d8c77f0/5B322027/t51.2885-19/s320x320/23668180_409038379513014_419175815813529600_n.jpg"
	caption := instagramMedia.Data.Caption.(map[string]interface{})

	data.Caption = caption["text"].(string)

	productsInShot := getProductsFromNames(watchesInshot)
	// if none is found, defualt to soho
	// if len(productsInShot) == 0 {
	//
	// }
	data.Products = productsInShot

	js, err := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

type instagramMediaData struct {
	PictureUrl                    string           `json:"picture_url"`
	WallaceHatchProfilePictureUrl string           `json:"wallace_hatch_profile_picture_url"`
	Location                      string           `json:"location"`
	Caption                       string           `json:"caption"`
	Products                      []stripe.Product `json:"products"`
}

func getInstagramMediaComments(mediaId string) (instagramCommentResp, error) {

	var msg instagramCommentResp
	client := &http.Client{}
	url := fmt.Sprint(instagramApiBaseURL, "media/", mediaId, "/comments?access_token=", instagramAccessToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Error with instagram url request", err)
		return msg, err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error with  executing instagram get media  request", err)
		return msg, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading all  body from instagram get media  ", err)
		return msg, err
	}

	err = json.Unmarshal(b, &msg)
	if err != nil {
		logger.Error("error unmarshaling instagram url shorten  data", err)
		return msg, err
	}

	return msg, err

}

func getInstagramMediaInfo(shortenUrl string) (instagramMediaResp, error) {
	var msg instagramMediaResp
	client := &http.Client{}
	url := fmt.Sprint(instagramApiBaseURL, "media/shortcode/", shortenUrl, "?access_token=", instagramAccessToken)
	logger.Info("url is : -> ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Error with instagram url request", err)
		return msg, err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error with  executing instagram get media  request", err)
		return msg, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading all  body from instagram get media  ", err)
		return msg, err
	}

	err = json.Unmarshal(b, &msg)
	if err != nil {
		logger.Error("error unmarshaling instagram url shorten  data", err)
		return msg, err
	}

	return msg, err

}
