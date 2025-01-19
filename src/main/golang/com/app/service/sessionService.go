package service

import (
	"github.com/mahbodkh/game-runner/model"
)

func GetSessionByPlayerId(userId int) model.Session {
	// probably we should check the userIdFrom or userIdTo
	return model.GetSessionByUserId(userId)
}

func GetSessionById(sessionId string) model.Session {
	return model.GetSessionById(sessionId)
}

func CreateSession(userId int32) model.Session {
	session := model.Session{UserIdFrom: userId}
	model.CreateSession(session)
	return session
}

func UpdateSession(session model.Session) {
	// validate the session if userIdFrom created and userIdTo is null
	// set the userIdTo to the session
	model.UpdateSession(session)
}

func validateSession(session model.Session) bool {
	// check this session never been created before
	// check the sessionId is not null and unique
	// probably alway userIdTo is null
	return model.ValidateSession(session)
}
