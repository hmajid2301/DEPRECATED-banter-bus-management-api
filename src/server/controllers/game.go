package controllers

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server/models"
)

// CreateGameType adds a new game type to the database.
func CreateGameType(_ *gin.Context, game *models.NewGame) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": game.Name,
	})
	gameLogger.Debug("Trying to add new game.")

	gameTag := getJSONTagFromStruct(game, "Name")

	var (
		emptyResponse struct{}
		filter        = map[string]string{gameTag: game.Name}
		item          *models.NewGame
	)

	err := database.Get("game", filter, &item)
	if err == nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game already exists.")

		return emptyResponse, errors.AlreadyExistsf("The game type %s", game.Name)
	}

	var EmptyGame = models.Game{
		Name:     game.Name,
		RulesURL: game.RulesURL,
		Questions: &models.Question{
			Rounds: &models.Rounds{
				One:   []string{},
				Two:   []string{},
				Three: []string{},
			},
		},
		Enabled: true,
	}

	err = database.Insert("game", EmptyGame)
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to add new game.")
	}
	return emptyResponse, nil
}

// GetAllGameType gets a list of names of all game types.
func GetAllGameType(_ *gin.Context) ([]string, error) {
	log.Debug("Trying to get all games.")
	games := []*models.Game{}
	err := database.GetAll("game", &games)

	if err != nil {
		log.Warn("Failed to get game types.")
	}

	var gameNames []string
	for _, game := range games {
		gameNames = append(gameNames, game.Name)
	}
	return gameNames, nil
}

// GetGameType gets all the information about a specific game type.
func GetGameType(_ *gin.Context, params *models.GameParams) (*models.Game, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Trying to get a game.")
	var game *models.Game
	gameTag := getJSONTagFromStruct(params, "Name")
	var filter = map[string]string{gameTag: params.Name}
	err := database.Get("game", filter, &game)
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game doesn't exists.")
		return game, errors.NotFoundf("The game type %s", params.Name)
	}
	return game, nil
}

// RemoveGameType removes a game type from the database.
func RemoveGameType(_ *gin.Context, params *models.GameParams) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameLogger.Debug("Removing game.")
	var emptyResponse struct{}
	gameTag := getJSONTagFromStruct(params, "Name")
	var filter = map[string]string{gameTag: params.Name}

	var game *models.Game
	err := database.Get("game", filter, &game)

	if err != nil {
		gameLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game doesn't exists.")
		return emptyResponse, errors.NotFoundf("The game type %s", params.Name)
	}

	database.Delete("game", filter)
	return emptyResponse, nil
}

// EnableGameType enables a disabled game type.
func EnableGameType(_ *gin.Context, params *models.GameParams) (struct{}, error) {
	log.WithFields(log.Fields{
		"game_name": params.Name,
	}).Info("Enabling game.")
	return updateGameType("true", params)
}

// DisableGameType disabled an enabled game type.
func DisableGameType(_ *gin.Context, params *models.GameParams) (struct{}, error) {
	log.WithFields(log.Fields{
		"game_name": params.Name,
	}).Info("Disabling game.")
	return updateGameType("false", params)
}

func updateGameType(enable string, params *models.GameParams) (struct{}, error) {
	gameLogger := log.WithFields(log.Fields{
		"game_name": params.Name,
	})
	gameTag := getJSONTagFromStruct(params, "Name")
	var filter = map[string]string{gameTag: params.Name}

	var game *models.Game
	err := database.Get("game", filter, &game)
	var emptyResponse struct{}
	if err != nil {
		gameLogger.WithFields(log.Fields{
			"game_name": params.Name,
		}).Warn("Game doesn't exist.")
		return emptyResponse, errors.NotFoundf("The game type %s", params.Name)
	} else if enable == strconv.FormatBool(game.Enabled) {
		gameLogger.WithFields(log.Fields{
			"game_name": params.Name,
		}).Warn("Game already exists.")
		return emptyResponse, errors.AlreadyExistsf("Game %s is already %s", params.Name, enable)
	}

	gameTagEnabled := getJSONTagFromStruct(game, "Enabled")
	b, err := strconv.ParseBool(enable)
	if err != nil {
		panic(err)
	}
	var update = map[string]bool{gameTagEnabled: b}

	updated, err := database.UpdateEntry("game", filter, update)
	if !updated || err != nil {
		gameLogger.Warn(
			fmt.Sprintf("Failed to update %s", game.Name),
		)
		return emptyResponse, errors.Errorf("Failed to update game with new question.")
	}

	return emptyResponse, nil
}

func getJSONTagFromStruct(model interface{}, fieldName string) string {
	field, ok := reflect.TypeOf(model).Elem().FieldByName(fieldName)
	if !ok {
		log.WithFields(log.Fields{
			"fieldName": field,
			"model":     model,
		}).Warn("Field not found.")
	}

	jsonFieldName := field.Tag.Get("json")
	return jsonFieldName
}
