package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// PrepareRequest creates the http request to be sent
func PrepareRequest(method, url, path, accessToken string, body *strings.Reader) (request *http.Request) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", url, path), body)
	if err != nil {
		panic("Failed while preparing a POST request")
	}
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("x-sentinel-tenant", "dev")
	request.Header.Add("Content-Type", "application/json")
	return
}

// ExecuteRequest calls the Do method on the request passed to it
func ExecuteRequest(request *http.Request) (response *http.Response) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic("Failed while performing Do on a request")
	}
	return
}

// ParseResponse unmarshals the response into the target interface
func ParseResponse(response *http.Response, expectedResponseStatus int) (data []byte) {
	data, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != expectedResponseStatus {
		panic(fmt.Sprintf(
			`Unexpected response code %s: %s`,
			strconv.Itoa(response.StatusCode),
			string(data)))
	}

	if err != nil {
		panic("Error encountered while parsing response " + err.Error())
	}

	return data
}
