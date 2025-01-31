package main

import (
	"fmt"

	"github.com/mahbodkh/game-runner/app/service"
)

func main() {
	fmt.Println("Starting the WebSocket server on :3000...")

	service.InitWebsocketService()
}
