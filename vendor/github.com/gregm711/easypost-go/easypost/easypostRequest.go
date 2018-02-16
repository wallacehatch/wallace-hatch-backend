package easypost

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const easypost_url = "https://api.easypost.com/v2"

var api_key string
var Request RequestControllerInterface

type RequestControllerInterface interface {
	do(method string, objectType string, objectUrl string, payload string) ([]byte, error)
}

func init() {
	Request = RequestControllerFake{}
}

func SetApiKey(key string) {
	api_key = key
}

type RequestController struct{}

//Request request an EasyPost API
func (rc RequestController) do(method string, objectType string, objectUrl string, payload string) ([]byte, error) {
	url := getRequestURL(objectType, objectUrl)
	body := strings.NewReader(payload)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New("Cannot create EasyPost request")
	}

	req.SetBasicAuth(api_key, "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("Cannot request EasyPost API")
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//getRequestUrl returns the correct url for EasyPost API
func getRequestURL(objectType string, objectUrl string) string {
	url := fmt.Sprintf("%vs", objectType)
	if objectType == "address" {
		url = "addresses"
	}
	if objectType == "batch" {
		url = "batches"
	}

	if objectUrl != "" {
		url = fmt.Sprintf("%v/%v", url, objectUrl)
	}

	return fmt.Sprintf("%v/%v", easypost_url, url)
}
