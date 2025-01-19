package main

import (
	"fmt"

	rest "github.com/mahbodkh/game-runner/api/rest"
	ws "github.com/mahbodkh/game-runner/api/ws"
)

func main() {
	fmt.Println("Starting the application...")

	ws.InitWebSocketService()
	rest.InitRestService()
}
