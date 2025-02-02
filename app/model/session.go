package model

/*
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50),
    current_player INT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/

import (
	"context"
	"time"

	"github.com/mahbodkh/game-runner/app/db"
)

const (
	StatusStarted    SessionStatus = "started"
	StatusInProgress SessionStatus = "in_progress"
	StatusFull       SessionStatus = "full"
	StatusInvalid    SessionStatus = "invalid"
	StatusError      SessionStatus = "error"
	StatusWaiting    SessionStatus = "waiting"
	StatusCreated    SessionStatus = "created"
	StatusMoved      SessionStatus = "moved"
	StatusCompleted  SessionStatus = "completed"
)

type SessionStatus string
type Session struct {
	ID            int32
	SessionID     int32
	GameId        int8 // 256 games should be enough for now
	UniqueId      int32
	CurrentPlayer int32
	Status        SessionStatus
	Players       []int32
	Board         [9]string
	Created       time.Time
	Updated       time.Time
}

func GetSessionById(id int32) (Session, error) {
	var session Session
	err := db.Conn.QueryRow(context.Background(), "SELECT id, current_player, status, created, updated FROM sessions WHERE id = $1", id).Scan(
		&session.ID, &session.CurrentPlayer, &session.Status, &session.Created, &session.Updated)
	if err != nil {
		return Session{}, err
	}

	// Get player IDs for the session
	rows, err := db.Conn.Query(context.Background(), "SELECT player_id FROM session_players WHERE session_id = $1", id)
	if err != nil {
		return Session{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var playerID int32
		if err := rows.Scan(&playerID); err != nil {
			return Session{}, err
		}
		session.Players = append(session.Players, playerID)
	}

	return session, nil
}

func SaveSession(session Session) (int32, error) {
	var id int32
	err := db.Conn.QueryRow(context.Background(),
		"INSERT INTO sessions (current_player, status, created, updated) VALUES ($1, $2, $3, $4) RETURNING id",
		session.CurrentPlayer, session.Status, session.Created, session.Updated).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Save player IDs for the session
	for _, playerID := range session.Players {
		_, err := db.Conn.Exec(context.Background(), "INSERT INTO session_players (session_id, player_id) VALUES ($1, $2)", id, playerID)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func UpdateSession(session Session) (Session, error) {
	_, err := db.Conn.Exec(context.Background(), "UPDATE sessions SET current_player=$1, status=$2, updated=$3 WHERE id=$4",
		session.CurrentPlayer, session.Status, session.Updated, session.ID)
	if err != nil {
		return Session{}, err
	}

	// Update player IDs for the session
	_, err = db.Conn.Exec(context.Background(), "DELETE FROM session_players WHERE session_id = $1", session.ID)
	if err != nil {
		return Session{}, err
	}

	for _, playerID := range session.Players {
		_, err := db.Conn.Exec(context.Background(), "INSERT INTO session_players (session_id, player_id) VALUES ($1, $2)", session.ID, playerID)
		if err != nil {
			return Session{}, err
		}
	}

	return session, nil
}

func DeleteSession(id int32) error {
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM sessions WHERE id = $1", id)
	return err
}

func GetSessionByUserId(userId int32) (Session, error) {
	var session Session
	err := db.Conn.QueryRow(context.Background(), "SELECT id, current_player, status, created, updated FROM sessions WHERE user_id = $1", userId).Scan(
		&session.ID, &session.CurrentPlayer, &session.Status, &session.Created, &session.Updated)
	if err != nil {
		return Session{}, err
	}

	// Get player IDs for the session
	rows, err := db.Conn.Query(context.Background(), "SELECT player_id FROM session_players WHERE session_id = $1", session.ID)
	if err != nil {
		return Session{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var playerID int32
		if err := rows.Scan(&playerID); err != nil {
			return Session{}, err
		}
		session.Players = append(session.Players, playerID)
	}

	return session, nil
}

func GetSessionByIdAndUserId(sessionId int32, userId int32) (Session, error) {
	var session Session
	err := db.Conn.QueryRow(context.Background(), "SELECT id, current_player, status, created, updated FROM sessions WHERE id = $1 AND user_id = $2", sessionId, userId).Scan(
		&session.ID, &session.CurrentPlayer, &session.Status, &session.Created, &session.Updated)
	if err != nil {
		return Session{}, err
	}

	// Get player IDs for the session
	rows, err := db.Conn.Query(context.Background(), "SELECT player_id FROM session_players WHERE session_id = $1", session.ID)
	if err != nil {
		return Session{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var playerID int32
		if err := rows.Scan(&playerID); err != nil {
			return Session{}, err
		}
		session.Players = append(session.Players, playerID)
	}

	return session, nil
}
