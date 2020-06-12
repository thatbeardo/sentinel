package handler

import (
	"errors"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	"github.com/gin-gonic/gin"
)

// NoRoute represents a 404 error for an invalid request URL
func NoRoute(c *gin.Context) {
	views.Wrap(errors.New("Path not found"), c)
}
