package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/factories"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// AddQuestion adds a new question to a game.
func (env *Env) AddQuestion(_ *gin.Context, questionInput *serverModels.AddQuestionInput) (string, error) {
	var (
		question = questionInput.NewQuestion
		gameName = questionInput.GameParams.Name
	)
	questionLogger := log.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": gameName,
	})
	questionLogger.Debug("Trying to add new question.")

	add := env.newGenericQuestion(question)
	err := env.validateQuestion(gameName, add)
	if err != nil {
		return "", err
	}

	q := service.QuestionService{
		DB:       env.DB,
		GameName: gameName,
		Question: add,
	}
	id, err := q.Add()

	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to add question.")
		return "", err
	}

	return id, nil
}

// RemoveQuestion removes a question  game.
func (env *Env) RemoveQuestion(_ *gin.Context, questionInput *serverModels.QuestionInput) error {
	var (
		questionID = questionInput.QuestionIDParams.ID
		gameName   = questionInput.GameParams.Name
	)
	questionLogger := log.WithFields(log.Fields{
		"question_id": questionID,
		"game_name":   gameName,
	})
	questionLogger.Debug("Trying to remove question.")

	q := service.QuestionService{
		DB:         env.DB,
		GameName:   gameName,
		QuestionID: questionID,
	}
	err := q.RemoveQuestion()
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to remove question.")
		return err
	}

	return nil
}

// AddTranslation adds a question in another language to a game.
func (env *Env) AddTranslation(_ *gin.Context, questionInput *serverModels.AddTranslationInput) error {
	var (
		questionID = questionInput.QuestionIDParams.ID
		question   = questionInput.QuestionTranslation
		gameName   = questionInput.GameParams.Name
		lang       = questionInput.LanguageParams.Language
	)

	questionLogger := log.WithFields(log.Fields{
		"question_id":   questionID,
		"game_name":     gameName,
		"language_code": lang,
		"new_question":  question.Content,
	})
	questionLogger.Debug("Trying to add new question translation.")

	_, err := language.Parse(lang)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err":           err,
			"language_code": lang,
		}).Warn("Bad language code.")
		return errors.BadRequestf("invalid language %s", lang)
	}

	q := service.QuestionService{
		DB:         env.DB,
		GameName:   gameName,
		QuestionID: questionID,
	}
	err = q.AddTranslation(question.Content, lang)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to add question.")
		return err
	}
	return nil
}

// RemoveTranslation removes a question in a language from a game.
func (env *Env) RemoveTranslation(_ *gin.Context, questionInput *serverModels.QuestionInput) error {
	var (
		questionID = questionInput.QuestionIDParams.ID
		gameName   = questionInput.GameParams.Name
		lang       = questionInput.LanguageParams.Language
	)
	questionLogger := log.WithFields(log.Fields{
		"question_id":   questionID,
		"game_name":     gameName,
		"language_code": lang,
	})
	questionLogger.Debug("Trying to remove question translation.")
	q := service.QuestionService{
		DB:         env.DB,
		GameName:   gameName,
		QuestionID: questionID,
	}
	err := q.RemoveTranslation(lang)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to remove question translation.")
		return err
	}

	return nil
}

// EnableQuestion enables a disabled game.
func (env *Env) EnableQuestion(_ *gin.Context, questionInput *serverModels.QuestionInput) (struct{}, error) {
	return env.updateEnable(questionInput, true)
}

// DisableQuestion disabled an enabled game.
func (env *Env) DisableQuestion(_ *gin.Context, questionInput *serverModels.QuestionInput) (struct{}, error) {
	return env.updateEnable(questionInput, false)
}

// GetAllGroups gets all group names from a certain round in a certain game
func (env *Env) GetAllGroups(_ *gin.Context, groupInput *serverModels.GroupInput) ([]string, error) {
	log.Debug("Trying to get all groups")
	q := service.QuestionService{
		DB:       env.DB,
		GameName: groupInput.Name,
	}
	groups, err := q.GetGroups(groupInput.Round)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get groups.")
		return []string{}, err
	}
	return groups, nil
}

func (env *Env) updateEnable(
	questionInput *serverModels.QuestionInput,
	enable bool,
) (struct{}, error) {
	var (
		questionID = questionInput.QuestionIDParams.ID
		gameName   = questionInput.GameParams.Name
	)
	questionLogger := log.WithFields(log.Fields{
		"question_id": questionID,
		"game_name":   gameName,
		"enable":      enable,
	})
	questionLogger.Debug("Trying to update question enable state.")

	var emptyResponse struct{}

	q := service.QuestionService{
		DB:         env.DB,
		GameName:   gameName,
		QuestionID: questionID,
	}
	updated, err := q.UpdateEnable(enable)
	if err != nil || !updated {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to update question state")
		return emptyResponse, err
	}

	return emptyResponse, nil
}

func (env *Env) newGenericQuestion(question serverModels.NewQuestion) models.GenericQuestion {
	group := &models.GenericQuestionGroup{
		Name: "",
		Type: "",
	}

	if question.Group != nil {
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

func (env *Env) validateQuestion(gameName string, question models.GenericQuestion) error {
	_, err := language.Parse(question.LanguageCode)
	if err != nil {
		return errors.BadRequestf("invalid language code %s", question.LanguageCode)
	}

	game, err := factories.GetGame(gameName)
	if err != nil {
		return err
	}

	err = game.ValidateQuestion(question)
	if err != nil {
		return err
	}

	return nil
}
