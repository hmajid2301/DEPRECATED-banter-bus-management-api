package controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/core"
	"banter-bus-server/src/server/models"
)

// AddQuestion adds a new question to a game.
func AddQuestion(_ *gin.Context, questionInput *models.ReceiveQuestionInput) error {
	var (
		question = questionInput.ReceiveQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question": question.Content,
	})
	questionLogger.Debug("Trying to add new question.")

	err := core.AddQuestion(game.Name, question.Round, question.Content)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to add question.")

		return err
	}

	return nil
}

// RemoveQuestion removes a question from a game.
func RemoveQuestion(_ *gin.Context, questionInput *models.ReceiveQuestionInput) error {
	var (
		question = questionInput.ReceiveQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question": question.Content,
	})
	questionLogger.Debug("Trying to remove question.")

	err := core.RemoveQuestion(game.Name, question.Round, question.Content)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to remove question.")

		return err
	}

	return nil
}
