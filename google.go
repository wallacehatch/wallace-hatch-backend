package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var googleApiKey string
var googleApiUrl string

const googleApiUrlBase = "https://www.googleapis.com/urlshortener/v1/url?key="

func init() {
	googleApiKey = os.Getenv("GOOGLE_URL_API")
	googleApiUrl = fmt.Sprint(googleApiUrlBase, googleApiKey)
}

type googleUrlResp struct {
	Kind    string `json:"kind"`
	Id      string `json:"id"`
	LongUrl string `json:"longUrl"`
	Error   struct {
		Errors []struct {
			Domain       string `json:"domain"`
			Reason       string `json:"reason"`
			Message      string `json:"message"`
			ExtendedHelp string `json:"extendedHelp"`
		} `json:"errors"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type googleUrlReq struct {
	LongUrl string `json:"longUrl"`
}

func shortenUrl(url string) (string, error) {
	var msg googleUrlResp
	data := googleUrlReq{LongUrl: url}
	b, err := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest("POST", googleApiUrl, bytes.NewBuffer(b))
	if err != nil {
		logger.Error("Error with google url shorten request", err)
		return msg.Id, err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error with  executing google url shorten  request", err)
		return msg.Id, err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading all  body from google url ", err)
	}

	err = json.Unmarshal(b, &msg)
	if err != nil {
		logger.Error("error unmarshaling google url shorten  data", err)
	}
	if msg.Error.Message != "" {
		logger.Error("Error from google", msg.Error)
		return msg.Id, errors.New(msg.Error.Message)
	}
	return msg.Id, err
}
