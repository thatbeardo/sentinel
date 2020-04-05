package permissions

import "github.com/gin-gonic/gin"

// PermissionRoutes sets up resource specific routes on the engine instance
func PermissionRoutes(r *gin.Engine) {
	router := r.Group("/permissions")
	router.GET("/", get())
}
