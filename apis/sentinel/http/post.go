package http

import (
	"io"
	"net/http"
)

// Post makes a request to the requested URL
func Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return http.Post(url, contentType, body)
}
