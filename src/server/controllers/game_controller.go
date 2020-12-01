package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/core"
	serverModels "banter-bus-server/src/server/models"
)

// CreateGame adds a new game.
func CreateGame(_ *gin.Context, game *serverModels.ReceiveGame) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": game.Name,
	})
	gameLogger.Debug("Trying to add new game.")

	var emptyResponse struct{}
	err := core.AddGame(game.Name, game.RulesURL)

	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to add the new game.")

		if errors.IsAlreadyExists(err) {
			gameLogger.WithFields(log.Fields{
				"err": err,
			}).Warn("Game already exists.")
		}
		return emptyResponse, err
	}

	return emptyResponse, nil
}

// GetAllGames gets a list of names of all game.
func GetAllGames(_ *gin.Context, params *serverModels.ListGameParams) ([]string, error) {
	log.Debug("Trying to get all games.")

	var (
		t = true
		f = false
	)

	var n *bool
	filters := map[string]*bool{
		"enabled":  &t,
		"disabled": &f,
		"all":      n,
	}

	enabledFilters := filters[params.Games]
	gameNames, err := core.GetAllGames(enabledFilters)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get game types.")
		return []string{}, err
	}

	return gameNames, nil
}

// GetGame gets all the information about a specific game.
func GetGame(_ *gin.Context, params *serverModels.GameParams) (*serverModels.Game, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Trying to get a game.")

	game, err := core.GetGame(params.Name)
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game doesn't exists.")
		return &serverModels.Game{}, errors.NotFoundf("The game %s", params.Name)
	}

	actualGame := &serverModels.Game{
		Name:     game.Name,
		RulesURL: game.RulesURL,
		Enabled:  *game.Enabled,
	}
	return actualGame, nil
}

// RemoveGame removes a game.
func RemoveGame(_ *gin.Context, params *serverModels.GameParams) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Removing game.")

	var emptyResponse struct{}

	err := core.RemoveGame(params.Name)
	if errors.IsNotFound(err) {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game doesn't exists.")

		return emptyResponse, err
	} else if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to remove game.")
		return emptyResponse, err
	}

	return emptyResponse, nil
}

// EnableGame enables a disabled game.
func EnableGame(_ *gin.Context, params *serverModels.GameParams) (struct{}, error) {
	return updateEnableGameState(params.Name, true)
}

// DisableGame disabled an enabled game.
func DisableGame(_ *gin.Context, params *serverModels.GameParams) (struct{}, error) {
	return updateEnableGameState(params.Name, false)
}

func updateEnableGameState(name string, enable bool) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": name,
		"enable":    enable,
	})
	log.WithFields(log.Fields{
		"game_name": name,
	}).Info("Trying to update enable state.")

	var emptyResponse struct{}

	updated, err := core.UpdateEnableGame(name, enable)
	if err != nil || !updated {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update game state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}
