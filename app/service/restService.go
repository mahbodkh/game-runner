package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type GameInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type PlayerInfo struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type CreateSessionRequest struct {
	Game       string `json:"game"`
	OpponentId int32  `json:"opponentId"`
}

func InitApiService() {
	http.HandleFunc("/api/games", handleListGames)
	http.HandleFunc("/api/players", handleListPlayers)
	http.HandleFunc("/api/sessions", handleCreateSession)

	logrus.Info("Apis started on http://localhost:3000/")
}

// GET /api/games
func handleListGames(w http.ResponseWriter, r *http.Request) {
	games := []GameInfo{
		{Name: "Tic-Tac-Toe 3x3", Path: "tic-toc-toe3x3"},
		{Name: "Tic-Tac-Toe 4x4", Path: "tic-toc-toe4x4"},
		{Name: "Tic-Tac-Toe 5x5", Path: "tic-toc-toe5x5"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

// GET /api/players?game=3x3
func handleListPlayers(w http.ResponseWriter, r *http.Request) {
	_ = r.URL.Query().Get("game")
	// Example: return all players or only those for the given game
	// We'll just return dummy data:
	players := []PlayerInfo{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

// POST /api/sessions  { "game":"3x3", "opponentId":2 }
func handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	// e.g. create or join a session in DB
	// For the example, we just say sessionId = 123
	sessionID := int32(123)
	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"sessionId":` + strconv.Itoa(int(sessionID)) + `}`))
}
