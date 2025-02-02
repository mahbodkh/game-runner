package main

import (
	"fmt"
	"net/http"

	"github.com/mahbodkh/game-runner/app/db"
	"github.com/mahbodkh/game-runner/app/service"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Starting the WebSocket server on :3000...")

	db.InitDB()

	service.InitWebsocketService()

	service.InitApiService()

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		logrus.Fatalf("ListenAndServe failed: %v", err)
	}
}
