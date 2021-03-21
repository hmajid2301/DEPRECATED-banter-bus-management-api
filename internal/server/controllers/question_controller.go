package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
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

	add := env.newGenericQuestion(question)
	q := service.QuestionService{
		Username: env.Conf.Official.Username,
		PoolName: env.Conf.Official.PoolName,
		DB:       env.DB,
		GameName: game.Name,
		Question: add,
	}
	err := q.Add()

	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to add question.")
		return err
	}

	return nil
}

// RemoveQuestion removes a question  game.
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

	remove := env.newGenericQuestion(question)
	q := service.QuestionService{
		DB:       env.DB,
		GameName: game.Name,
		Question: remove,
		Username: env.Conf.Official.Username,
		PoolName: env.Conf.Official.PoolName,
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
		question = questionInput.QuestionTranslation
		game     = questionInput.GameParams
		lang     = questionInput.LanguageParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":      question.OriginalQuestion.Content,
		"game_name":     game.Name,
		"language_code": lang.Language,
	})
	questionLogger.Debug("Trying to add new question translation.")

	newLanguage := lang.Language
	_, err := language.Parse(newLanguage)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err":           err,
			"language_code": newLanguage,
		}).Warn("Bad language code.")
		return errors.BadRequestf("invalid language %s", newLanguage)
	}

	genericQuestion := env.newGenericQuestion(question.OriginalQuestion)
	q := service.QuestionService{
		Username: env.Conf.Official.Username,
		PoolName: env.Conf.Official.PoolName,
		DB:       env.DB,
		GameName: game.Name,
		Question: genericQuestion,
	}
	err = q.AddTranslation(question.NewQuestion.Content, newLanguage)
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
		question = questionInput.NewQuestion
		game     = questionInput.GameParams
		lang     = questionInput.LanguageParams
	)
	questionLogger := log.WithFields(log.Fields{
		"question":      question.Content,
		"game_name":     game.Name,
		"language_code": lang.Language,
	})
	questionLogger.Debug("Trying to remove question translation.")

	question.LanguageCode = lang.Language
	remove := env.newGenericQuestion(question)
	q := service.QuestionService{
		DB:       env.DB,
		GameName: game.Name,
		Question: remove,
		Username: env.Conf.Official.Username,
		PoolName: env.Conf.Official.PoolName,
	}
	err := q.RemoveTranslation()
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
		Username: env.Conf.Official.Username,
		PoolName: env.Conf.Official.PoolName,
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

	update := env.newGenericQuestion(question)
	q := service.QuestionService{
		DB:       env.DB,
		GameName: game.Name,
		Question: update,
		Username: env.Conf.Official.Username,
		PoolName: env.Conf.Official.PoolName,
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
