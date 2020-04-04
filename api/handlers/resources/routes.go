package resources

import "github.com/gin-gonic/gin"

// ResourceRoutes sets up resource specific routes on the engine instance
func ResourceRoutes(r *gin.Engine) {
	router := r.Group("/resources")
	router.GET("/", get())
}
