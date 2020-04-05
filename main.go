package main

import (
	"github.com/thatbeardo/go-sentinel/api/handlers/server"
	"github.com/thatbeardo/go-sentinel/pkg/resource"
)

func main() {
	shutdown, session := server.Initialize()
	defer shutdown()

	resourceRepository := resource.NewNeo4jRepository(session)
	resourceService := resource.NewService(resourceRepository)

	// neo4jrepository := resources

	engine := server.SetupRouter(resourceService)
	server.Orchestrate(engine, shutdown)
}
