package main

import (
	"github.com/thatbeardo/go-sentinel/pkg/resource"
	"github.com/thatbeardo/go-sentinel/server"
)

func main() {
	shutdown, session := server.Initialize()

	resourceRepository := resource.NewNeo4jRepository(session)
	resourceService := resource.NewService(resourceRepository)

	engine := server.SetupRouter(resourceService)
	server.Orchestrate(engine, shutdown)
}
