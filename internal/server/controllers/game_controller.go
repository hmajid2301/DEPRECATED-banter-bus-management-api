package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service"
)

// CreateGame adds a new game.
func (env *Env) CreateGame(_ *gin.Context, game *serverModels.NewGame) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"name": game.Name,
	})
	gameLogger.Debug("Trying to add new game.")

	var emptyResponse struct{}
	g := service.GameService{DB: env.DB, Name: game.Name}
	err := g.Add(game.RulesURL)

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

	enabled := filters[params.Games]
	g := service.GameService{DB: env.DB}
	names, err := g.GetAll(enabled)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get game.")
		return []string{}, err
	}

	return names, nil
}

// GetGame gets all the information about a specific game.
func (env *Env) GetGame(_ *gin.Context, params *serverModels.GameParams) (*serverModels.Game, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Trying to get a game.")

	g := service.GameService{DB: env.DB, Name: params.Name}
	game, err := g.Get()
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game does not exist.")
		return &serverModels.Game{}, errors.NotFoundf("The game %s", params.Name)
	}

	gameObj := &serverModels.Game{
		Name:     game.Name,
		RulesURL: game.RulesURL,
		Enabled:  *game.Enabled,
	}
	return gameObj, nil
}

// RemoveGame removes a game.
func (env *Env) RemoveGame(_ *gin.Context, params *serverModels.GameParams) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Removing game.")

	var emptyResponse struct{}

	gameService := service.GameService{DB: env.DB, Name: params.Name}
	err := gameService.Remove()
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
	gameService := service.GameService{DB: env.DB, Name: name}
	updated, err := gameService.UpdateEnable(enable)

	if err != nil || !updated {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update game state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}
