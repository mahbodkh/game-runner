package model

import (
	"time"
)

type Game struct {
	ID          int32     `json:"id" yaml:"id" boil:"id" db:"id" toml:"id"`
	Name        string    `json:"name" yaml:"name" boil:"name" db:"name" toml:"name"`
	Description string    `json:"description" yaml:"description" boil:"description" db:"description" toml:"description"`
	Created     time.Time `json:"created" yaml:"created" boil:"created" db:"created" toml:"created"`
	Updated     time.Time `json:"updated" yaml:"updated" boil:"updated" db:"updated" toml:"updated"`
}

var GameColumns = struct {
	ID          string
	Name        string
	Description string
	Created     string
	Updated     string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
	Created:     "created",
	Updated:     "updated",
}

var GameTableColumns = struct {
	ID          string
	Name        string
	Description string
	Created     string
	Updated     string
}{
	ID:          "game.id",
	Name:        "game.name",
	Description: "game.description",
	Created:     "game.created",
	Updated:     "game.updated",
}

func GetGameList() []Game {
	return []Game{}
}

func GetGameById(id string) Game {
	return Game{}
}
