package views

import (
	"net/http"
	"strings"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/gin-gonic/gin"
)

// ErrView represents the error structure sent as a response
type ErrView struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Source Source `json:"source"`
	Detail string `json:"detail"`
}

// Source provides details about where the error started from
type Source struct {
	Pointer   string `json:"pointer"`
	Parameter string `json:"parameter"`
}

// ErrHTTPStatusMap determines the status code to be sent back to the user based on the error message
var ErrHTTPStatusMap = map[string]int{
	models.ErrNoContent.Error():    http.StatusNotFound,
	models.ErrDatabase.Error():     http.StatusInternalServerError,
	models.ErrNotFound.Error():     http.StatusNotFound,
	models.ErrPathNotFound.Error(): http.StatusNotFound,
}

// Wrap generates Errors to be returned to the user
func Wrap(err error, c *gin.Context) {
	msg := err.Error()
	code := ErrHTTPStatusMap[msg]
	if code == 0 {
		if strings.Contains(msg, "Key:") && strings.Contains(msg, "Error:") {
			code = http.StatusBadRequest
		} else {
			code = http.StatusInternalServerError
		}
	}
	c.AbortWithStatusJSON(code, GenerateErrorResponse(code, msg, c.FullPath()))
}

// GenerateErrorResponse creates an error payload to be sent to the user
func GenerateErrorResponse(code int, msg string, path string) ErrView {
	source := Source{
		Pointer:   path,
		Parameter: "query-parameter-todo",
	}
	return ErrView{
		ID:     "error-id-todo",
		Status: code,
		Source: source,
		Detail: msg,
	}
}
