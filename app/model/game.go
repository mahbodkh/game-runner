package model

import (
	"time"
)

type Game struct {
	ID          int32
	Name        string
	Description string
	Path        string
	Created     time.Time
	Updated     time.Time
}
