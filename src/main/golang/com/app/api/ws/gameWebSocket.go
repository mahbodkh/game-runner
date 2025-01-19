package api

import (
	"github.com/mahbodkh/game-runner/service"
)

func InitWebSocketService() {
	return "Websocket service is running..."
}

// initGame - make a session for the starter player and create a unique link for seccond player to join
func InitWebSocket() {
	// check the player user
	// check the initilizer player for the game
	// create a session for the player
	// create a unique link for the seccond player to join
	// init websocket for the game
	service.InitGameService()
}
