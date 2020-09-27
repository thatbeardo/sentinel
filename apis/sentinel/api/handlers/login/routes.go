package login

import (
	"github.com/gin-gonic/gin"
)

// Routes sets up context specific routes on the engine instance
func Routes(r *gin.RouterGroup) {
	r.POST("/login", getAccessToken())
}
