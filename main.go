package main

import "github.com/thatbeardo/go-sentinel/api/handlers/server"

func main() {
	engine := server.SetupRouter()
	engine.Run()
}
