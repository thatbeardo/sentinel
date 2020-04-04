package main

import "github.com/thatbeardo/go-sentinel/api/handlers/server"

func main() {
	server.ConnectToDB()
	engine := server.SetupRouter()
	engine.Run()
}
