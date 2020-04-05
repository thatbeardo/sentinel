package main

import (
	"github.com/thatbeardo/go-sentinel/api/handlers/server"
	"github.com/thatbeardo/go-sentinel/pkg/resource"
)

func main() {
	shutdown, session := server.Initialize()

	resourceRepository := resource.NewNeo4jRepository(session)
	resourceService := resource.NewService(resourceRepository)

	engine := server.SetupRouter(resourceService)
	server.Orchestrate(engine, shutdown)
}
