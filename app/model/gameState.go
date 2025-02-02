package model

import (
	"time"
)

type Action string

type GameState struct {
	Board        [9]string     `json:"board"`
	Player       int32         `json:"player"`
	Action       Action        `json:"action"`
	SessionID    int32         `json:"sessionId,omitempty"`
	Index        string        `json:"index,omitempty"`
	Status       SessionStatus `json:"status"`
	IsWin        bool          `json:"isWin"`
	IsPlayerTurn bool          `json:"isPlayerTurn"`
	Message      string        `json:"message"`
	Created      time.Time     `json:"created"`
	Updated      time.Time     `json:"updated"`
}

const (
	ActionGetSessions = "getSessions"
	ActionJoin        = "join"
	ActionMove        = "move"
)
