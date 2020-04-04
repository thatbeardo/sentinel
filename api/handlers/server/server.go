package server

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"
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
