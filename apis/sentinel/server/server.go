package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bithippie/guard-my-app/apis/sentinel/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/alexcesaro/statsd.v2"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// GenerateRouter instantiates and initializes a new Router.
func GenerateRouter(r *gin.Engine) *gin.RouterGroup {
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "x-sentinel-tenant", "Authorization")
	config.AllowMethods = append(config.AllowMethods, "OPTIONS")
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	setupSwagger(r)
	return r.Group("/v1")
}

// GenerateStatsdClient instantiates and returns a statsd client
func GenerateStatsdClient(host, port string) (*statsd.Client, error) {
	return statsd.New(
		statsd.Address(
			fmt.Sprintf("%s:%s", host, port)))
}

func setupSwagger(r *gin.Engine) {
	hostURL := fmt.Sprintf("https://%s", os.Getenv("HOST"))
	docs.SwaggerInfo.Host = hostURL
	r.StaticFile("/docs", "./docs/swagger.json")

	url := ginSwagger.URL(fmt.Sprintf("%s/docs", hostURL))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// Orchestrate begins listening on PORT and gracefully shuts down the server incase of interrupt
func Orchestrate(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
