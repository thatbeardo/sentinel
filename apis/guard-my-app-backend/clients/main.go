package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/patrickmn/go-cache"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var AUTH0_CLIENT_ID = os.Getenv("AUTH0_CLIENT_ID")
var AUTH0_CLIENT_SECRET = os.Getenv("AUTH0_CLIENT_SECRET")
var AUTH0_DOMAIN = os.Getenv("AUTH0_DOMAIN")
var SENTINEL_API_AUDIENCE = os.Getenv("SENTINEL_API_AUDIENCE")

func getM2MToken() (jwt Auth0Token, err error) {
	url := fmt.Sprintf(`https://%s/oauth/token`, AUTH0_DOMAIN)

	payload := strings.NewReader(
		fmt.Sprintf(
			`grant_type=client_credentials&client_id=%s&client_secret=%s&audience=%s`,
			AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, SENTINEL_API_AUDIENCE,
		),
	)

	fmt.Println(payload)

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var auth0Token = Auth0Token{}
	err = json.Unmarshal(body, &auth0Token)
	jwt = auth0Token
	return
}

type Auth0Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	jwt, err := getM2MToken()

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": jwt.AccessToken,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	// validate jwt
	// create application in auth0
	// create client grant with no scopes to sentinel api
	// call sentinel api to create resource client
	// call sentinel api to create resource org
	// create resource tenant owned by org
	// create policy
	// grant policy to client
	// permit policy to perform admin operations on organization
	// call auth0 to make a new client id and client secret
	// call sentinel api with client id NOT secret
	lambda.Start(Handler)
}
