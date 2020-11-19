package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/core"
	"banter-bus-server/src/server/models"
)

// CreateGameType adds a new game type to the database.
func CreateGameType(_ *gin.Context, game *models.NewGame) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": game.Name,
	})
	gameLogger.Debug("Trying to add new game type.")

	var emptyResponse struct{}
	err := core.AddGameType(game.Name, game.RulesURL)

	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to add new game.")

		if errors.IsAlreadyExists(err) {
			gameLogger.WithFields(log.Fields{
				"err": err,
			}).Warn("Game already exists.")
		}
		return emptyResponse, err
	}

	return emptyResponse, nil
}

// GetAllGameTypes gets a list of names of all game types.
func GetAllGameTypes(_ *gin.Context) ([]string, error) {
	log.Debug("Trying to get all games.")

	gameNames, err := core.GetAllGameTypes()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get game types.")
		return []string{}, err
	}

	return gameNames, nil
}

// GetGameType gets all the information about a specific game type.
func GetGameType(_ *gin.Context, params *models.GameParams) (*models.Game, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Trying to get a game.")

	game, err := core.GetGameType(params.Name)
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game doesn't exists.")
		return &models.Game{}, errors.NotFoundf("The game type %s", params.Name)
	}

	actualGame := &models.Game{
		Name:     game.Name,
		RulesURL: game.RulesURL,
		Enabled:  *game.Enabled,
	}
	return actualGame, nil
}

// RemoveGameType removes a game type from the database.
func RemoveGameType(_ *gin.Context, params *models.GameParams) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Removing game.")

	var emptyResponse struct{}

	err := core.RemoveGameTypes(params.Name)
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

// EnableGameType enables a disabled game type.
func EnableGameType(_ *gin.Context, params *models.GameParams) (struct{}, error) {
	return updateEnableGameState(params.Name, true)
}

// DisableGameType disabled an enabled game type.
func DisableGameType(_ *gin.Context, params *models.GameParams) (struct{}, error) {
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

	updated, err := core.UpdateEnableGameType(name, enable)
	if err != nil || !updated {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update game state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}
