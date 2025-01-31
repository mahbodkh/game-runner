package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for simplicity. Restrict in production.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}
		log.Printf("Received: %s", message)

		response := fmt.Sprintf("Thanks: %s", message)
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}

	log.Println("Client disconnected")
}

func InitWebsocketService() {

	http.HandleFunc("/ws", handleWS)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/html/tic-toc-toe"))
	http.Handle("/", fs)

	fmt.Println("WebSocket service started on", "ws://localhost:3000/ws")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}
