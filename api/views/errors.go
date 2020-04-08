package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	models "github.com/thatbeardo/go-sentinel/models"
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
	models.ErrNoContent.Error(): http.StatusNotFound,
	models.ErrDatabase.Error():  http.StatusInternalServerError,
}

// Wrap generates Errors to be returned to the user
func Wrap(err error, c *gin.Context) {
	msg := err.Error()
	code := ErrHTTPStatusMap[msg]
	if code == 0 {
		code = http.StatusInternalServerError
	}
	c.AbortWithStatusJSON(code, generateErrorResponse(code, msg, c.FullPath()))
}

func generateErrorResponse(code int, msg string, path string) ErrView {
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
