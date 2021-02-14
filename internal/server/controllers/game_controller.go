package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// CreateGame adds a new game.
func (env *Env) CreateGame(_ *gin.Context, game *serverModels.NewGame) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": game.Name,
	})
	gameLogger.Debug("Trying to add new game.")

	var emptyResponse struct{}

	gameService := biz.GameService{DB: env.DB}
	err := gameService.Add(game.Name, game.RulesURL)

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
func (env *Env) GetAllGames(_ *gin.Context, params *serverModels.ListGameParams) ([]string, error) {
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
	gameService := biz.GameService{DB: env.DB}
	gameNames, err := gameService.GetAll(enabledFilters)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get game.")
		return []string{}, err
	}

	return gameNames, nil
}

// GetGame gets all the information about a specific game.
func (env *Env) GetGame(_ *gin.Context, params *serverModels.GameParams) (*serverModels.Game, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Trying to get a game.")

	gameService := biz.GameService{DB: env.DB}
	game, err := gameService.Get(params.Name)
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game does not exist.")
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
func (env *Env) RemoveGame(_ *gin.Context, params *serverModels.GameParams) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Removing game.")

	var emptyResponse struct{}

	gameService := biz.GameService{DB: env.DB}
	err := gameService.Remove(params.Name)
	if errors.IsNotFound(err) {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game does not exist.")

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
func (env *Env) EnableGame(_ *gin.Context, params *serverModels.GameParams) (struct{}, error) {
	return env.updateEnableGameState(params.Name, true)
}

// DisableGame disabled an enabled game.
func (env *Env) DisableGame(_ *gin.Context, params *serverModels.GameParams) (struct{}, error) {
	return env.updateEnableGameState(params.Name, false)
}

func (env *Env) updateEnableGameState(name string, enable bool) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": name,
		"enable":    enable,
	})
	log.Debug("Trying to update enable state.")

	var emptyResponse struct{}

	gameService := biz.GameService{DB: env.DB}
	updated, err := gameService.UpdateEnable(name, enable)
	if err != nil || !updated {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update game state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}
