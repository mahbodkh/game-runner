package api

import (
	"github.com/mahbodkh/game-runner/service"
)

func InitRestService() {
	return "Rest service is running..."
}

// initGame - make a session for the starter player and create a unique link for seccond player to join
func InitGame() {
	// check the player user
	// check the initilizer player for the game
	// create a session for the player
	// create a unique link for the seccond player to join
	// init websocket for the game
	service.InitGameService()
}

func JoinGame(sessionId string) {
	// check the player user
	// check the session
	// join the game by return websocket connection to the seccond player
	service.StartGameService()
}

func LookingForOnlinePlayers() {
	// check the player user
	// list of online players
	service.StartGame()
}
