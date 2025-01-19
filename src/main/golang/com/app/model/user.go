package model

import (
	"time"
)

type User struct {
	ID           int32
	TelegramId   int32 // not sure if this is the right type
	Name         string
	LanguageCode string
	Created      time.Time
	Updated      time.Time
}
