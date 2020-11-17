package core

import (
	"fmt"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/dbmodels"

	"github.com/juju/errors"
)

// AddQuestion is add questions to a game type.
func AddQuestion(gameName string, round string, content string) error {
	game, _ := GetGameType(gameName)
	if game.Name == "" {
		return errors.NotFoundf("The game type %s", gameName)
	} else if !*game.Enabled {
		return errors.AlreadyExistsf("The game type %s is not enabled", gameName)
	}

	questionExists := doesQuestionExist(gameName, round, content)
	if questionExists {
		return errors.AlreadyExistsf("The question already exists")
	}

	filter := dbmodels.Game{Name: gameName}
	questionToAdd := map[string]string{fmt.Sprintf("questions.rounds.%s", round): content}
	updated, err := database.AppendToEntry("game", filter, questionToAdd)
	if !updated || err != nil {
		return errors.Errorf("Failed to update game with new question.")
	}

	return nil
}

func doesQuestionExist(gameName string, round string, content string) bool {
	var (
		filter = map[string]string{"name": gameName, fmt.Sprintf("questions.rounds.%s", round): content}
		game   *dbmodels.Game
	)

	err := database.Get("game", filter, &game)
	if err != nil || game.Name == "" {
		return false
	}

	return true
}
