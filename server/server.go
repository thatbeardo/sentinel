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

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/thatbeardo/go-sentinel/handlers/resources"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

// SetupRouter instantiates and initializes a new Router.
func SetupRouter(service resource.Service) *gin.Engine {
	r := gin.Default()

	resources.ResourceRoutes(r, service)

	return r
}

// Initialize connects to the database and returns a shut down function
func Initialize() (func(), neo4j.Session) {
	session, driver, err := ConnectToDB()
	fmt.Println(err)
	return func() {
		session.Close()
		driver.Close()
	}, session
}

// ConnectToDB establishes connection to the neo4j database
func ConnectToDB() (neo4j.Session, neo4j.Driver, error) {
	// define driver, session and result vars
	var (
		driver  neo4j.Driver
		session neo4j.Session
		err     error
	)
	// initialize driver to connect to localhost with ID and password
	if driver, err = neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "password", "")); err != nil {
		return nil, nil, err
	}
	// Open a new session with write access
	if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
		return nil, nil, err
	}
	return session, driver, nil
}

// Orchestrate begins listening on 8080 and gracefully shuts down the server incase of interrupt
func Orchestrate(router *gin.Engine, cleanup func()) {
	srv := &http.Server{
		Addr:    ":8080",
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
	defer cleanup()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
