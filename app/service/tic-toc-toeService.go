package service

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mahbodkh/game-runner/app/model"
	"github.com/sirupsen/logrus"
)

const (
	X = "X"
	O = "O"
)

func AddPlayer(conn *websocket.Conn, name string) (model.Player, error) {
	logrus.Debugf("Adding player %s", name)
	player := model.Player{
		Name:       name,
		Mark:       X,
		TelegramId: 0,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	// Add player to the database
	if player.Name == "" {
		logrus.Error("Invalid player data provided: Name is missing!")
		return model.Player{}, errors.New("invalid player data provided: Name is missing")
	}

	playerID, err := model.SavePlayer(player)
	if err != nil {
		return model.Player{}, err
	}
	player.ID = playerID

	players[conn] = player

	// Check for an existing session with only one player
	var session *model.Session
	for _, s := range sessions {
		if s.Status == model.StatusWaiting && len(s.Players) == 1 {
			session = &s
			logrus.Debugf("Found session %d for new player %d", s.ID, playerID)
			break
		}
	}

	if session == nil {
		// Create a new session for the player
		session = &model.Session{
			CurrentPlayer: player.ID,
			Status:        model.StatusWaiting,
			Players:       []int32{player.ID},
			Created:       time.Now(),
			Updated:       time.Now(),
		}
		// Add session to the database
		sessionID, err := model.SaveSession(*session)
		if err != nil {
			return model.Player{}, err
		}
		session.ID = sessionID
		sessions[session.ID] = *session
	} else {
		// Add the player to the existing session
		player.Mark = O
		session.Players = append(session.Players, player.ID)
		session.Status = model.StatusInProgress
		session.CurrentPlayer = session.Players[0]
		session.Updated = time.Now()
		sessions[session.ID] = *session
		// Update the session in the database ??!
		_, err := model.UpdateSession(*session)
		if err != nil {
			logrus.Error("Error updating session:", err)
			return model.Player{}, err
		}

		// Notify both players that the game has started
		broadcastGameStart(session.ID)
	}

	conn.WriteJSON(model.GameState{
		Status:       model.StatusCreated,
		SessionID:    session.ID,
		Player:       player.ID,
		IsPlayerTurn: true,
	})

	return player, nil
}

func handleMakeMove(conn *websocket.Conn, player model.Player, gameState *model.GameState) {
	session := getSessionFromCacheOrDbByUserId(gameState.SessionID, player.ID)

	if len(session.Players) < 2 {
		conn.WriteJSON(model.GameState{
			Status:  model.StatusError,
			Message: "This session doesn't have two players yet!",
		})
		return
	}

	// its wrong to check for the current player here because the client doesn't know who's turn it is

	if session.CurrentPlayer != player.ID {
		conn.WriteJSON(model.GameState{
			Status:  model.StatusError,
			Message: "Not your turn",
		})
		return
	}

	index, err := strconv.Atoi(gameState.Index)
	if err != nil {
		log.Println("Invalid move index:", err)
		return
	}

	if gameState.Board[index] == "" {
		gameState.Board[index] = player.Mark
		if len(session.Players) == 2 {
			if player.Mark == X {
				// index out of range [1] with length 1
				// something is wrong here with the logic because it's always X and client doesn't know!
				session.CurrentPlayer = session.Players[1]
			} else {
				session.CurrentPlayer = session.Players[0]
			}
		}

		if checkWin(gameState.Board, player.Mark) {
			session.Status = model.StatusCompleted
			broadcastGameResult(gameState.SessionID, player.ID, "win")
		} else if checkDraw(gameState.Board) {
			session.Status = model.StatusCompleted
			broadcastGameResult(gameState.SessionID, 0, "draw")
		}

		sessions[gameState.SessionID] = session
		broadcastBoardUpdate(gameState.SessionID, gameState.Board)
	}
}

func getSessionFromCacheOrDbByUserId(sessionId int32, playerId int32) model.Session {
	session, ok := sessions[sessionId]
	if !ok {
		var err error
		if playerId != 0 {
			session, err = model.GetSessionByIdAndUserId(sessionId, playerId)
		} else {
			session, err = model.GetSessionById(sessionId)
		}
		if err != nil {
			logrus.Error("Error getting session:", err)
			return model.Session{}
		}
		sessions[session.ID] = session
	}
	return session
}

func getSessionFromCacheOrDbById(sessionId int32) model.Session {
	session, ok := sessions[sessionId]
	if !ok {
		var err error
		sessionDb, err := model.GetSessionById(sessionId)
		if err != nil {
			logrus.Error("Error getting session:", err)
			return model.Session{}
		}
		session = sessionDb
		sessions[session.ID] = session
	}
	return session
}

func checkWin(board [9]string, mark string) bool {
	winPatterns := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, pattern := range winPatterns {
		if board[pattern[0]] == mark && board[pattern[1]] == mark && board[pattern[2]] == mark {
			return true
		}
	}
	return false
}

func checkDraw(board [9]string) bool {
	for _, cell := range board {
		if cell == "" {
			return false
		}
	}
	return true
}

func broadcastGameResult(sessionID int32, winnerID int32, result string) {
	for conn, player := range players {
		if contains(sessions[sessionID].Players, player.ID) {
			gameResult := map[string]interface{}{
				"status":  model.StatusCompleted,
				"winner":  winnerID,
				"result":  result,
				"session": sessionID,
			}
			conn.WriteJSON(gameResult)
		}
	}
}

func broadcastBoardUpdate(sessionID int32, board [9]string) {
	for conn, player := range players {
		if contains(sessions[sessionID].Players, player.ID) {
			boardUpdate := map[string]interface{}{
				"status":  model.StatusMoved,
				"board":   board,
				"session": sessionID,
			}
			conn.WriteJSON(boardUpdate)
		}
	}
}
