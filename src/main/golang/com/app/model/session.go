package model

import (
	"database/sql"
	"time"
)

var db *sql.DB

type Session struct {
	ID         int32
	UserIdFrom int32
	UserIdTo   int32
	GameId     int8 // 256 games should be enough for now
	UniqueId   int32

	Created time.Time
	Updated time.Time
}

func GetSessionByUserId(userId int) Session {

	return Session{}
}

func GetSessionById(sessionId string) Session {

	return Session{}
}

func CreateSession(session Session) {

}

func UpdateSession(session Session) {

}

func ValidateSession(session Session) bool {

	return true
}
