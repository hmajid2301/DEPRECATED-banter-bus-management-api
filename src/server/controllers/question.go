package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server/models"
)

// AddQuestion adds a new question to a game.
func AddQuestion(_ *gin.Context, questionInput *models.NewQuestionInput) (struct{}, error) {
	var (
		question = questionInput.NewQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question": question.Content,
	})
	questionLogger.Debug("Trying to add new question.")

	questionTag := getJSONTagFromStruct(&models.Game{}, "Questions")
	roundTag := getJSONTagFromStruct(&models.Question{}, "Rounds")

	var (
		emptyResponse    struct{}
		questionLocation = fmt.Sprintf("%s.%s.%s", questionTag, roundTag, question.Round)
		questionToAdd    = map[string]string{questionLocation: question.Content}
		gameName         = models.GameParams{Name: game.Name}
		currentGame      *models.Game
	)

	err := database.Get("game", gameName, &currentGame)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game doesn't exists.")

		return emptyResponse, errors.NotFoundf("The game type %s", game.Name)
	}

	if !currentGame.Enabled {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Game isn't enabled.")

		return emptyResponse, errors.AlreadyExistsf("The game type %s not enabled", game.Name)
	}

	err = database.Get("game", questionToAdd, &models.Game{})
	if err == nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Question already exists.")

		return emptyResponse, errors.AlreadyExistsf(
			"The question %s in game type %s, and round %s", question.Content, game.Name, question.Round,
		)
	}

	updated, err := database.AppendToEntry("game", gameName, questionToAdd)
	if !updated || err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Error(
			fmt.Sprintf("Failed to update %s, with new question %s", questionInput.GameParams.Name, questionInput.NewQuestion),
		)

		return emptyResponse, errors.Errorf("Failed to update game with new question.")
	}
	return emptyResponse, nil
}
