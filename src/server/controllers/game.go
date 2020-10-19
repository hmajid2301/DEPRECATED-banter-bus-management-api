package controllers

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server/models"
)

func CreateGameType(c *gin.Context, game *models.NewGame) (struct{}, error) {
	var emptyResponse struct{}

	gameTag := getJSONTagFromStruct(game, "Name")
	var filter = map[string]string{gameTag: game.Name}
	var item *models.NewGame

	err := database.Get("game", filter, &item)
	if err == nil {
		return emptyResponse, errors.AlreadyExistsf("The game type %s", game.Name)
	}

	var EmptyGame = models.Game{
		Name:      game.Name,
		RulesURL:  game.RulesURL,
		Questions: []models.Question{},
		Enabled:   true,
	}

	err = database.Insert("game", EmptyGame)
	if err != nil {
		fmt.Println(err)
	}
	return emptyResponse, nil
}

func GetAllGameType(c *gin.Context) ([]string, error) {
	games := []*models.Game{}
	database.GetAll("game", &games)
	var gameNames []string
	for _, game := range games {
		gameNames = append(gameNames, game.Name)
	}
	return gameNames, nil
}

func GetGameType(c *gin.Context, params *models.GameParams) (*models.Game, error) {
	var game *models.Game
	gameTag := getJSONTagFromStruct(params, "Name")
	var filter = map[string]string{gameTag: params.Name}
	err := database.Get("game", filter, &game)
	if err != nil {
		return game, errors.NotFoundf("The game type %s", params.Name)
	}
	return game, nil
}

func RemoveGameType(c *gin.Context, params *models.GameParams) (struct{}, error) {
	var emptyResponse struct{}
	gameTag := getJSONTagFromStruct(params, "Name")
	var filter = map[string]string{gameTag: params.Name}

	var game *models.Game
	err := database.Get("game", filter, &game)
	if err != nil {
		return emptyResponse, errors.NotFoundf("The game type %s", params.Name)
	}
	database.Delete("game", filter)
	return emptyResponse, nil
}

func EnableGameType(c *gin.Context, params *models.GameParams) (struct{}, error) {
	return updateGameType(true, params)
}

func DisableGameType(c *gin.Context, params *models.GameParams) (struct{}, error) {
	return updateGameType(false, params)
}

func updateGameType(enable bool, params *models.GameParams) (struct{}, error) {
	gameTag := getJSONTagFromStruct(params, "Name")
	var filter = map[string]string{gameTag: params.Name}

	var game *models.Game
	err := database.Get("game", filter, &game)
	var emptyResponse struct{}
	if err != nil {
		return emptyResponse, errors.NotFoundf("The game type %s", params.Name)
	} else if enable == game.Enabled {
		enabledString := "enabled"
		if !enable {
			enabledString = "disabled"
		}
		return emptyResponse, errors.NewAlreadyExists(errors.New(""), fmt.Sprintf("Game %s is already %s", params.Name, enabledString))
	}

	gameTagEnabled := getJSONTagFromStruct(game, "Enabled")
	var update = map[string]bool{gameTagEnabled: enable}

	database.PartialUpdate("game", filter, update)
	return emptyResponse, nil

}

func getJSONTagFromStruct(model interface{}, fieldName string) string {
	field, ok := reflect.TypeOf(model).Elem().FieldByName(fieldName)
	if !ok {
		panic("Field not found")
	}
	jsonFieldName := string(field.Tag.Get("json"))
	return jsonFieldName
}
