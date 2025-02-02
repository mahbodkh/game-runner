package model

/*
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    telegram_id INT,
    name VARCHAR(255),
    language_code VARCHAR(10),
    session_id INT,
    mark CHAR(1),
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/

import (
	"context"
	"time"

	"github.com/mahbodkh/game-runner/app/db"
)

type Player struct {
	ID           int32
	TelegramId   int32
	Name         string
	LanguageCode string
	Mark         string
	Created      time.Time
	Updated      time.Time
}

func GetPlayerById(id int32) (Player, error) {
	var player Player
	err := db.Conn.QueryRow(context.Background(), "SELECT id, telegram_id, name, language_code, mark, created, updated FROM players WHERE id = $1", id).Scan(
		&player.ID, &player.TelegramId, &player.Name, &player.LanguageCode, &player.Mark, &player.Created, &player.Updated)
	if err != nil {
		return Player{}, err
	}
	return player, nil
}

func SavePlayer(player Player) (int32, error) {
	var id int32
	err := db.Conn.QueryRow(context.Background(),
		"INSERT INTO players (telegram_id, name, language_code, mark, created, updated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		player.TelegramId, player.Name, player.LanguageCode, player.Mark, player.Created, player.Updated).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdatePlayer(player Player) (Player, error) {
	_, err := db.Conn.Exec(context.Background(), "UPDATE players SET telegram_id=$1, name=$2, language_code=$3, mark=$4, updated=$5 WHERE id=$6",
		player.TelegramId, player.Name, player.LanguageCode, player.Mark, player.Updated, player.ID)
	if err != nil {
		return Player{}, err
	}
	return player, nil
}

func DeletePlayer(id int32) error {
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM players WHERE id = $1", id)
	return err
}
