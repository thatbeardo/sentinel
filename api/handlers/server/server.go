package server

import (
	"github.com/gin-gonic/gin"
	"github.com/thatbeardo/go-sentinel/api/handlers/permissions"
	"github.com/thatbeardo/go-sentinel/api/handlers/resources"
)

// SetupRouter instantiates and initializes a new Router.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	resources.ResourceRoutes(r)
	permissions.PermissionRoutes(r)

	return r
}
