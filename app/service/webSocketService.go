package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mahbodkh/game-runner/app/model"
	"github.com/sirupsen/logrus"
)

var sessions = make(map[int32]model.Session)
var players = make(map[*websocket.Conn]model.Player)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Upgrade error")
		return
	}
	defer conn.Close()
	logrus.Info("Client connected")
	handleClient(conn)
}

func handleClient(conn *websocket.Conn) {
	// Read the player's name
	_, message, err := conn.ReadMessage()
	if err != nil {
		logrus.WithError(err).Error("Read message error")
		conn.Close()
		return
	}
	playerName := string(message)
	// Add the player
	player, err := AddPlayer(conn, playerName)
	if err != nil {
		logrus.WithError(err).Error("Error adding player")
		conn.Close()
		return
	}
	logrus.Infof("Player %s added with ID %d", player.Name, player.ID)

	// Main loop to read messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logrus.WithError(err).Error("Read message error")
			conn.Close()
			break
		}

		if messageType == websocket.TextMessage {
			logrus.Infof("Received message: %s", message)
			logrus.Debugf("Player: %s", player.Name)
			var gameState model.GameState
			err := json.Unmarshal(message, &gameState)
			if err != nil {
				logrus.WithError(err).Error("Error parsing game state")
				continue
			}

			logrus.Infof("Game state action: %s", gameState.Action)

			switch gameState.Action {
			case model.ActionGetSessions:
				sendAvailableSessions(conn)
			case model.ActionJoin:
				handleJoinSession(conn, player, gameState.SessionID)
			case model.ActionMove:
				handleMakeMove(conn, player, &gameState)
			default:
				logrus.Warn("Unknown action: ", gameState.Action)
			}
		}
	}
}

func sendAvailableSessions(conn *websocket.Conn) {
	var availableSessions []model.Session
	for _, session := range sessions {
		if session.Status == model.StatusWaiting {
			playersInSession := []int32{}
			for _, playerID := range session.Players {
				playersInSession = append(playersInSession, playerID)
			}
			availableSessions = append(availableSessions, model.Session{
				ID:        session.ID,
				Status:    session.Status,
				Players:   playersInSession,
				Created:   session.Created,
				Updated:   session.Updated,
				GameId:    session.GameId,
				UniqueId:  session.UniqueId,
				SessionID: session.SessionID,
			})
		}
	}
	conn.WriteJSON(availableSessions)
}

func handleJoinSession(conn *websocket.Conn, player model.Player, sessionId int32) {
	session := getSessionFromCacheOrDbByUserId(sessionId, player.ID)
	if session.ID == 0 {
		conn.WriteJSON(model.GameState{
			Status: model.StatusInvalid,
		})
		return
	}

	if len(session.Players) == 1 {
		player.Mark = O
	}

	if session.Status == model.StatusWaiting {
		session.Players = append(session.Players, player.ID)
		session.Status = model.StatusInProgress
		session.CurrentPlayer = session.Players[0]
		session.Updated = time.Now()
		sessions[session.ID] = session

		logrus.Infof("Player %s joined session %d", player.Name, sessionId)

		// Notify both players that the game has started
		broadcastGameStart(session.ID)
	}

	conn.WriteJSON(model.GameState{
		Status:       model.StatusCreated,
		SessionID:    session.ID,
		Player:       player.ID,
		IsPlayerTurn: true,
	})
}

func broadcastGameStart(sessionId int32) {
	session, exists := sessions[sessionId]
	if !exists {
		logrus.Warnf("Session %d does not exist", sessionId)
		return
	}

	for conn, player := range players {
		if contains(session.Players, player.ID) {
			gameState := model.GameState{
				Status:       model.StatusStarted,
				SessionID:    sessionId,
				IsPlayerTurn: player.ID == session.CurrentPlayer,
			}
			conn.WriteJSON(gameState)
		}
	}
}

func contains(slice []int32, item int32) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func InitWebsocketService() {
	http.HandleFunc("/ws", handleWS)

	fs := http.FileServer(http.Dir("./web/html/"))
	http.Handle("/", fs)

	logrus.Info("WebSocket service started on ws://localhost:3000/ws")
}
