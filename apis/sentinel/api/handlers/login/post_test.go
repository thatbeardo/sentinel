package login_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

const validRequest = "{\"client_id\":\"client-id\", \"client_secret\":\"client-secret\"}"

var testErr = errors.New("some-error")

func TestGetAccessToken_BadRequestSent_ReturnsStatusBadRequest(t *testing.T) {
	router := setupRouter()
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/login", "{}")
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			http.StatusBadRequest,
			"Key: 'ClientCredentials.ClientID' Error:Field validation for 'ClientID' failed on the 'required' tag\nKey: 'ClientCredentials.ClientSecret' Error:Field validation for 'ClientSecret' failed on the 'required' tag",
			"/v1/login"),
		http.StatusBadRequest)
}

func TestGetAccessToken_MarshallingFails_ReturnsStatusInternalServerError(t *testing.T) {
	defer injection.Reset()
	injection.Marshal = func(v interface{}) ([]byte, error) {
		return nil, testErr
	}

	router := setupRouter()
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/login", validRequest)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			http.StatusInternalServerError,
			"some-error",
			"/v1/login"),
		http.StatusInternalServerError)
}

func TestGetAccessToken_PostRequestFails_ReturnInternalServerError(t *testing.T) {
	defer injection.Reset()
	injection.Post = func(string, string, io.Reader) (*http.Response, error) {
		return nil, testErr
	}

	router := setupRouter()
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/login", validRequest)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			http.StatusInternalServerError,
			"some-error",
			"/v1/login"),
		http.StatusInternalServerError)
}

func TestGetAccessToken_ReadAllFails_ReturnInternalServerError(t *testing.T) {
	defer injection.Reset()
	injection.Unmarshal = func([]byte, interface{}) error {
		return testErr
	}

	router := setupRouter()
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/login", validRequest)
	defer cleanup()

	testutil.ValidateResponse(
		t,
		response,
		views.GenerateErrorResponse(
			http.StatusInternalServerError,
			"some-error",
			"/v1/login"),
		http.StatusInternalServerError)
}

func TestGetAccessToken_AllOK_ReturnStatusOK(t *testing.T) {
	router := setupRouter()
	response, cleanup := testutil.PerformRequest(router, "POST", "/v1/login", validRequest)
	defer cleanup()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
