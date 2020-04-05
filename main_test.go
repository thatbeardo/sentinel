// package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/thatbeardo/go-sentinel/api/handlers/server"
// )

// func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
// 	req, _ := http.NewRequest(method, path, nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)
// 	return w
// }

// func TestResourcesOk(t *testing.T) {
// 	// Build our expected body
// 	body := gin.H{
// 		"name": "Harshil",
// 	}
// 	// Grab our router
// 	router := server.SetupRouter()
// 	// Perform a GET request with that handler.
// 	w := performRequest(router, "GET", "/resources/")
// 	// Assert we encoded correctly,
// 	// the request gives a 200
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	// Convert the JSON response to a map
// 	var response map[string]string
// 	err := json.Unmarshal([]byte(w.Body.String()), &response)
// 	// Grab the value & whether or not it exists
// 	value, exists := response["name"]
// 	// Make some assertions on the correctness of the response.
// 	assert.Nil(t, err)
// 	assert.True(t, exists)
// 	assert.Equal(t, body["name"], value)
// }

// func TestPermissionsOk(t *testing.T) {
// 	// Build our expected body
// 	body := gin.H{
// 		"message": "permissions",
// 	}
// 	// Grab our router
// 	router := server.SetupRouter()
// 	// Perform a GET request with that handler.
// 	w := performRequest(router, "GET", "/permissions/")
// 	// Assert we encoded correctly,
// 	// the request gives a 200
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	// Convert the JSON response to a map
// 	var response map[string]string
// 	err := json.Unmarshal([]byte(w.Body.String()), &response)
// 	// Grab the value & whether or not it exists
// 	value, exists := response["message"]
// 	// Make some assertions on the correctness of the response.
// 	assert.Nil(t, err)
// 	assert.True(t, exists)
// 	assert.Equal(t, body["message"], value)
// }
