package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// AddQuestion adds a new question to a game.
func (env *Env) AddQuestion(_ *gin.Context, questionInput *serverModels.QuestionInput) error {
	var (
		question = questionInput.NewQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": game.Name,
	})
	questionLogger.Debug("Trying to add new question.")

	questionToAdd := env.newGenericQuestion(question)
	questionService := biz.QuestionService{DB: env.DB}
	err := questionService.Add(game.Name, questionToAdd)

	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to add question.")
		return err
	}

	return nil
}

// UpdateQuestion adds a new question to a game.
func (env *Env) UpdateQuestion(_ *gin.Context, questionInput *serverModels.UpdateQuestionInput) error {
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
		questionLogger.WithFields(log.Fields{
			"err":           err,
			"language_code": newQuestionLanguage,
		}).Warn("Bad language code.")
		return errors.BadRequestf("Invalid language %s", newQuestionLanguage)
	}

	existingQuestion := env.newGenericQuestion(question.OriginalQuestion)
	questionService := biz.QuestionService{DB: env.DB}
	err = questionService.Update(
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
func (env *Env) RemoveQuestion(_ *gin.Context, questionInput *serverModels.QuestionInput) error {
	var (
		question = questionInput.NewQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": game.Name,
	})
	questionLogger.Debug("Trying to remove question.")

	questionToRemove := env.newGenericQuestion(question)
	questionService := biz.QuestionService{DB: env.DB}
	err := questionService.Remove(game.Name, questionToRemove)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to remove question.")
		return err
	}

	return nil
}

// EnableQuestion enables a disabled game.
func (env *Env) EnableQuestion(_ *gin.Context, questionInput *serverModels.QuestionInput) (struct{}, error) {
	return env.updateEnableQuestionState(questionInput, true)
}

// DisableQuestion disabled an enabled game.
func (env *Env) DisableQuestion(_ *gin.Context, questionInput *serverModels.QuestionInput) (struct{}, error) {
	return env.updateEnableQuestionState(questionInput, false)
}

// GetAllGroups gets all group names from a certain round in a certain game
func (env *Env) GetAllGroups(_ *gin.Context, groupInput *serverModels.GroupInput) ([]string, error) {
	log.Debug("Trying to get all groups")
	questionService := biz.QuestionService{DB: env.DB}
	groups, err := questionService.GetGroups(groupInput.GameName, groupInput.Round)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get groups.")
		return []string{}, err
	}
	return groups, nil
}

func (env *Env) updateEnableQuestionState(
	questionInput *serverModels.QuestionInput,
	enable bool,
) (struct{}, error) {
	var (
		question = questionInput.NewQuestion
		game     = questionInput.GameParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": game.Name,
		"enable":    enable,
	})
	questionLogger.Debug("Trying to update question enable state.")

	var emptyResponse struct{}

	questionToUpdate := env.newGenericQuestion(question)
	questionService := biz.QuestionService{DB: env.DB}
	updated, err := questionService.UpdateEnable(game.Name, enable, questionToUpdate)
	if err != nil || !updated {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update question state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}

func (env *Env) newGenericQuestion(question serverModels.NewQuestion) models.GenericQuestion {
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
