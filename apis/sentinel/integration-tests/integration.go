package main

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/integration-tests/admin"
	"fmt"
	"runtime/debug"
	"os"
)

func main() {
	defer handlePanicIfAny()
	admin.ScaffoldNewClient("client", "environment")

}

func handlePanicIfAny() {
	if r := recover(); r != nil {
		fmt.Println(r)
		fmt.Println(string(debug.Stack()))
		os.Exit(1)
	}
}