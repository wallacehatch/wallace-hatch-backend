package main

import (
	"encoding/json"
	"fmt"
	"github.com/Vorkytaka/instagram-go-scraper/instagram"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var matchTags = regexp.MustCompile(`[^\S]|^#([^\s#.,!)]+)$`)

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
	caption, ok := instagramMedia.Data.Caption.(map[string]interface{})
	if !ok {
		logger.Error("Error with instagram media, using scraping method ")
		instagramMedia = scrapeInstagramPost(shortenedUrl)
		caption = instagramMedia.Data.Caption.(map[string]interface{})

	}
	watchesInshot := make([]string, 0)
	for _, tag := range instagramMedia.Data.Tags {

		if watchTags[tag.(string)] == true {
			watchesInshot = append(watchesInshot, tag.(string))
		}
	}

	data.PictureUrl = instagramMedia.Data.Images.StandardResolution.URL
	data.WallaceHatchProfilePictureUrl = "https://instagram.fcmh1-1.fna.fbcdn.net/vp/0e5a33e58f4e512e19b038747d8c77f0/5B322027/t51.2885-19/s320x320/23668180_409038379513014_419175815813529600_n.jpg"

	data.Caption = caption["text"].(string)

	loc, ok := instagramMedia.Data.Location.(map[string]interface{})
	if ok {
		data.Location = loc["name"].(string)
	}

	productsInShot := getProductsFromNames(watchesInshot)

	// if none is found, defualt to kallio rose
	if len(productsInShot) == 0 {
		productsInShot = getProductsFromNames([]string{"kalliorose"})
	}

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

func scrapeInstagramPost(mediaId string) instagramMediaResp {
	instagramData := instagramMediaResp{}
	url := fmt.Sprint("https://www.instagram.com/p/", mediaId, "/")
	media, _ := instagram.GetMediaByURL(url)
	tags := getTags(media.Caption)
	captionMap := make(map[string]interface{})
	captionMap["text"] = media.Caption
	instagramData.Data.Caption = captionMap
	new := make([]interface{}, len(tags))
	for i, v := range tags {
		new[i] = v
	}
	instagramData.Data.Tags = new
	instagramData.Data.Images.StandardResolution.URL = media.MediaList[0].URL
	return instagramData
}

func getTags(s string) []string {
	res := make([]string, 0)
	fields := strings.FieldsFunc(s, tagsSplitter)
	for _, v := range fields {
		sub := matchTags.FindStringSubmatch(v)
		if len(sub) > 1 {
			res = append(res, strings.ToLower(sub[1]))
		}
	}
	return res
}

func tagsSplitter(c rune) bool {
	if unicode.IsSpace(c) {
		return true
	}
	switch c {
	case '.', ',', '!', ')':
		return true
	}
	return false
}
