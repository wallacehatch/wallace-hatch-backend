package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

	fmt.Println(watchesInshot)

	// comments, err := getInstagramMediaComments(instagramMedia.Data.ID)
	// logger.Info("Comments!!!! ", comments)

	js, err := json.Marshal(instagramMedia)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

type instagramMediaData struct {
	PictureUrl                    string
	wallaceHatchProfilePictureUrl string
	location                      string
	caption                       string
}

func getInstagramMediaComments(mediaId string) (instagramCommentResp, error) {

	var msg instagramCommentResp
	client := &http.Client{}
	url := fmt.Sprint(instagramApiBaseURL, "media/", mediaId, "/comments?access_token=", instagramAccessToken)
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
