package service

import (
	"github.com/mahbodkh/game-runner/model"
)

func GetGameList() []model.Game {
	return model.GetGameList()
}

func GetGameById(id string) model.Game {
	return model.GetGameById(id)
}

func InitGameService() {

}
