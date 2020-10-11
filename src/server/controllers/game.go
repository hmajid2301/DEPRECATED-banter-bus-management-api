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
	var item models.NewGame

	err := database.Get("game", filter, item)
	if err == nil {
		return emptyResponse, errors.AlreadyExistsf("The game type %s", game.Name)
	}

	var EmptyGame = models.Game{
		Name:      game.Name,
		RulesURL:  game.RulesURL,
		Questions: []models.Question{},
	}

	err = database.Insert("game", EmptyGame)
	if err != nil {
		fmt.Println(err)
	}
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
