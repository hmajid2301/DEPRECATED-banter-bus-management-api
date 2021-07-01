package questions

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// QuestionAPI is related to all the data the controllers need.
type QuestionAPI struct {
	Conf   core.Conf
	Logger *log.Logger
	DB     database.Database
}

// AddQuestion adds a new question to a game.
func (env *QuestionAPI) AddQuestion(_ *gin.Context, questionInput *AddQuestionInput) (string, error) {
	var (
		question = questionInput.QuestionIn
		gameName = questionInput.GameParams.Name
	)
	questionLogger := env.Logger.WithFields(log.Fields{
		"question":  question.Content,
		"game_name": gameName,
	})
	questionLogger.Debug("Trying to add new question.")

	err := env.validateQuestion(gameName, question)
	add := env.newGenericQuestion(question)

	if err != nil {
		return "", err
	}

	q := QuestionService{
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

func (env *QuestionAPI) validateQuestion(gameName string, question QuestionIn) error {
	languageCode := question.LanguageCode
	if languageCode == "" {
		languageCode = "en"
	}

	_, err := language.Parse(languageCode)
	if err != nil {
		return errors.BadRequestf("invalid language code %s", languageCode)
	}

	game, err := GetGame(gameName)
	if err != nil {
		return err
	}

	err = game.ValidateQuestion(question)
	if err != nil {
		return err
	}

	return nil
}

func (env *QuestionAPI) newGenericQuestion(question QuestionIn) GenericQuestion {
	group := &GenericQuestionGroup{
		Name: "",
		Type: "",
	}

	if question.Group != nil {
		group = &GenericQuestionGroup{
			Name: question.Group.Name,
			Type: question.Group.Type,
		}
	}

	if question.LanguageCode == "" {
		question.LanguageCode = "en"
	}

	newQuestion := GenericQuestion{
		Content:      question.Content,
		Round:        question.Round,
		Group:        group,
		LanguageCode: question.LanguageCode,
	}

	return newQuestion
}

// RemoveQuestion removes a question  game.
func (env *QuestionAPI) RemoveQuestion(_ *gin.Context, questionInput *QuestionInput) error {
	var (
		questionID = questionInput.QuestionIDParams.ID
		gameName   = questionInput.GameParams.Name
	)
	questionLogger := env.Logger.WithFields(log.Fields{
		"question_id": questionID,
		"game_name":   gameName,
	})
	questionLogger.Debug("Trying to remove question.")

	q := QuestionService{
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
func (env *QuestionAPI) AddTranslation(_ *gin.Context, questionInput *AddTranslationInput) error {
	var (
		questionID = questionInput.QuestionIDParams.ID
		question   = questionInput.QuestionTranslationIn
		gameName   = questionInput.GameParams.Name
		lang       = questionInput.LanguageParams.Language
	)

	questionLogger := env.Logger.WithFields(log.Fields{
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

	q := QuestionService{
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
func (env *QuestionAPI) RemoveTranslation(_ *gin.Context, questionInput *QuestionInput) error {
	var (
		questionID = questionInput.QuestionIDParams.ID
		gameName   = questionInput.GameParams.Name
		lang       = questionInput.LanguageParams.Language
	)
	questionLogger := env.Logger.WithFields(log.Fields{
		"question_id":   questionID,
		"game_name":     gameName,
		"language_code": lang,
	})
	questionLogger.Debug("Trying to remove question translation.")
	q := QuestionService{
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
func (env *QuestionAPI) EnableQuestion(_ *gin.Context, questionInput *QuestionInput) (struct{}, error) {
	return env.updateEnable(questionInput, true)
}

// DisableQuestion disabled an enabled game.
func (env *QuestionAPI) DisableQuestion(_ *gin.Context, questionInput *QuestionInput) (struct{}, error) {
	return env.updateEnable(questionInput, false)
}

// GetAllGroups gets all group names from a certain round in a certain game
func (env *QuestionAPI) GetAllGroups(_ *gin.Context, groupInput *GroupInput) ([]string, error) {
	gameName := groupInput.GameParams.Name
	questionLogger := env.Logger.WithFields(log.Fields{
		"game_name": gameName,
		"round":     groupInput.Round,
	})

	questionLogger.Debug("Trying to get all groups")
	q := QuestionService{
		DB:       env.DB,
		GameName: gameName,
	}
	groups, err := q.GetGroups(groupInput.Round)

	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get groups.")
		return []string{}, err
	}
	return groups, nil
}

func (env *QuestionAPI) updateEnable(
	questionInput *QuestionInput,
	enable bool,
) (struct{}, error) {
	var (
		questionID = questionInput.QuestionIDParams.ID
		gameName   = questionInput.GameParams.Name
	)
	questionLogger := env.Logger.WithFields(log.Fields{
		"question_id": questionID,
		"game_name":   gameName,
		"enable":      enable,
	})
	questionLogger.Debug("Trying to update question enable state.")

	var emptyResponse struct{}

	q := QuestionService{
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
