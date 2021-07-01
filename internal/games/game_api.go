package games

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"gitlab.com/banter-bus/banter-bus-management-api/internal"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

type GameAPI struct {
	Conf   core.Conf
	Logger *log.Logger
	DB     database.Database
}

func (env *GameAPI) AddGame(_ *gin.Context, game *GameIn) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"name": game.Name,
	})
	gameLogger.Debug("Trying to add new game.")

	var emptyResponse struct{}
	g := GameService{DB: env.DB, Name: game.Name}
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

func (env *GameAPI) GetGames(_ *gin.Context, params *ListGameParams) ([]string, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"games": params.Games,
	})
	gameLogger.Debug("Trying to get all games.")

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
	g := GameService{DB: env.DB}
	names, err := g.GetAll(enabled)
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get game.")
		return []string{}, err
	}

	return names, nil
}

func (env *GameAPI) GetGame(_ *gin.Context, params *internal.GameParams) (*GameOut, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Trying to get a game.")

	g := GameService{DB: env.DB, Name: params.Name}
	game, err := g.Get()
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game does not exist.")
		return &GameOut{}, errors.NotFoundf("The game %s", params.Name)
	}

	gameObj := &GameOut{
		Name:     game.Name,
		RulesURL: game.RulesURL,
		Enabled:  *game.Enabled,
	}
	return gameObj, nil
}

func (env *GameAPI) RemoveGame(_ *gin.Context, params *internal.GameParams) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Removing game.")

	var emptyResponse struct{}

	gameService := GameService{DB: env.DB, Name: params.Name}
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

func (env *GameAPI) EnableGame(_ *gin.Context, params *internal.GameParams) (struct{}, error) {
	return env.updateEnableGameState(params.Name, true)
}

func (env *GameAPI) DisableGame(_ *gin.Context, params *internal.GameParams) (struct{}, error) {
	return env.updateEnableGameState(params.Name, false)
}

func (env *GameAPI) updateEnableGameState(name string, enable bool) (struct{}, error) {
	gameLogger := env.Logger.WithFields(log.Fields{
		"game_name": name,
		"enable":    enable,
	})
	gameLogger.Debug("Trying to update enable state.")

	var emptyResponse struct{}
	gameService := GameService{DB: env.DB, Name: name}
	updated, err := gameService.UpdateEnable(enable)

	if err != nil || !updated {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update game state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}
