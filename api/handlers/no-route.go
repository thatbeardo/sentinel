package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/api/views"
)

// NoRoute represents a 404 error for an invalid request URL
func NoRoute(c *gin.Context) {
	views.Wrap(errors.New("Path not found"), c)
	return
}
