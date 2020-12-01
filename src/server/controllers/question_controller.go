package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"banter-bus-server/src/core"
	"banter-bus-server/src/core/models"
	serverModels "banter-bus-server/src/server/models"
)

// TODO: Add New Question Content -> EN, default EN
// TODO: Update with translation
// TODO: Filter using Structs

// AddQuestion adds a new question to a game.
func AddQuestion(_ *gin.Context, questionInput *serverModels.ReceiveQuestionInput) error {
	var (
		question = questionInput.ReceiveQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": game.Name,
	})
	questionLogger.Debug("Trying to add new question.")

	questionToAdd := newGenericQuestion(question)
	err := core.AddQuestion(game.Name, questionToAdd)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to add question.")
		return err
	}

	return nil
}

// UpdateQuestion adds a new question to a game.
func UpdateQuestion(_ *gin.Context, questionInput *serverModels.UpdateQuestionInput) error {
	var (
		question = questionInput.QuestionTranslation
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.OriginalQuestion.Content,
		"game_name": game.Name,
	})
	questionLogger.Debug("Trying to add new question.")

	newQuestionLanguage := question.NewQuestion.LanguageCode
	_, err := language.Parse(newQuestionLanguage)
	if err != nil {
		return errors.BadRequestf("Invalid language %s", newQuestionLanguage)
	}

	existingQuestion := newGenericQuestion(question.OriginalQuestion)
	err = core.UpdateQuestion(
		game.Name,
		existingQuestion,
		question.NewQuestion.Content,
		newQuestionLanguage,
	)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to add question.")
		return err
	}

	return nil
}

// RemoveQuestion removes a question from a game.
func RemoveQuestion(_ *gin.Context, questionInput *serverModels.ReceiveQuestionInput) error {
	var (
		question = questionInput.ReceiveQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": game.Name,
	})
	questionLogger.Debug("Trying to remove question.")

	questionToRemove := newGenericQuestion(question)
	err := core.RemoveQuestion(game.Name, questionToRemove)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to remove question.")

		return err
	}

	return nil
}

// EnableQuestion enables a disabled game.
func EnableQuestion(_ *gin.Context, questionInput *serverModels.ReceiveQuestionInput) (struct{}, error) {
	return updateEnableQuestionState(questionInput, true)
}

// DisableQuestion disabled an enabled game.
func DisableQuestion(_ *gin.Context, questionInput *serverModels.ReceiveQuestionInput) (struct{}, error) {
	return updateEnableQuestionState(questionInput, false)
}

func updateEnableQuestionState(questionInput *serverModels.ReceiveQuestionInput, enable bool) (struct{}, error) {
	var (
		question = questionInput.ReceiveQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": game.Name,
		"enable":    enable,
	})
	questionLogger.Debug("Trying to update question enable state.")

	var emptyResponse struct{}

	questionToUpdate := newGenericQuestion(question)
	updated, err := core.UpdateEnableQuestion(game.Name, enable, questionToUpdate)
	if err != nil || !updated {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update question state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}

func newGenericQuestion(question serverModels.ReceiveQuestion) models.GenericQuestion {
	var group *models.GenericQuestionGroup
	if question.Group == nil {
		group = &models.GenericQuestionGroup{
			Name: "",
			Type: "",
		}
	} else {
		group = &models.GenericQuestionGroup{
			Name: question.Group.Name,
			Type: question.Group.Type,
		}
	}

	if question.LanguageCode == "" {
		question.LanguageCode = "en"
	}

	newQuestion := models.GenericQuestion{
		Content:      question.Content,
		Round:        question.Round,
		Group:        group,
		LanguageCode: question.LanguageCode,
	}

	return newQuestion
}
