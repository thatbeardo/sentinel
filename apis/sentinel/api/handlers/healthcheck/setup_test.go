package healthcheck_test

import (
	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/healthcheck"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	testutil.RemoveMiddleware()
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	healthcheck.Routes(group)
	return r
}
