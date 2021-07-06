package questions

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"gitlab.com/banter-bus/banter-bus-management-api/internal"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

type QuestionAPI struct {
	Conf   core.Conf
	Logger *log.Logger
	DB     database.Database
}

func (env *QuestionAPI) AddQuestion(_ *gin.Context, questionInput *AddQuestionInput) (string, error) {
	var (
		question = questionInput.QuestionIn
		gameName = questionInput.Name
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
	err := q.Remove()
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to remove question.")
		return err
	}

	return nil
}

func (env *QuestionAPI) GetQuestion(_ *gin.Context, questionInput *GetQuestionInput) (QuestionGenericOut, error) {
	var (
		questionID   = questionInput.ID
		gameName     = questionInput.Name
		languageCode = questionInput.Language
	)
	questionLogger := env.Logger.WithFields(log.Fields{
		"question_id":   questionID,
		"game_name":     gameName,
		"language_code": languageCode,
	})
	questionLogger.Debug("Trying to get question.")

	_, err := language.Parse(languageCode)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err":           err,
			"language_code": languageCode,
		}).Warn("Bad language code.")
		return QuestionGenericOut{}, errors.BadRequestf("invalid language %s", languageCode)
	}

	q := QuestionService{
		DB:         env.DB,
		GameName:   gameName,
		QuestionID: questionID,
	}
	question, err := q.Get()
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to get question.")
		return QuestionGenericOut{}, err
	}

	questionOut, err := newGenericQuestionOut(question, languageCode)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to get question.")
	}
	return questionOut, err
}

func newGenericQuestionOut(question Question, languageCode string) (QuestionGenericOut, error) {
	content := question.Content[languageCode]
	if content == "" {
		return QuestionGenericOut{}, errors.NotFoundf("language code %s not found for question with id %s", languageCode, question.ID)
	}

	questionOut := QuestionGenericOut{
		Content: content,
		Round:   question.Round,
		Enabled: *question.Enabled,
	}

	if question.Group != nil {
		questionOut.Group = &QuestionGroupInOut{
			Name: question.Group.Name,
			Type: question.Group.Type,
		}
	}

	return questionOut, nil
}

func (env *QuestionAPI) GetQuestions(_ *gin.Context, params *ListQuestionParams) ([]QuestionOut, error) {
	questionLogger := env.Logger.WithFields(log.Fields{
		"game_name":     params.Name,
		"round":         params.Round,
		"language_code": params.Language,
		"group_name":    params.GroupName,
		"random":        params.Random,
		"enabled":       params.Enabled,
		"limit":         params.Limit,
	})

	// todo refactor this
	questionLogger.Debug("Trying to get questions.")

	_, err := language.Parse(params.Language)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err":           err,
			"language_code": params.Language,
		}).Warn("Bad language code.")
		return []QuestionOut{}, errors.BadRequestf("invalid language %s", params.Language)
	}

	q := QuestionService{
		DB:       env.DB,
		GameName: params.Name,
	}

	enabled := internal.GetEnabledBool(params.Enabled)
	searchParams := SearchParams{
		Round:     params.Round,
		GroupName: params.GroupName,
		Random:    params.Random,
		Enabled:   enabled,
		Limit:     params.Limit,
		Language:  params.Language,
	}

	questions, err := q.GetList(searchParams)
	if err != nil {
		questionLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to get questions.")
		return []QuestionOut{}, err
	}

	questionOut := newQuestionOut(questions, params.Language)
	return questionOut, nil
}

func newQuestionOut(questions Questions, languageCode string) []QuestionOut {
	questionsOut := []QuestionOut{}

	for _, question := range questions {
		questionType := "question"
		if question.Round == "answers" || question.Round == "drawing" {
			questionType = "answer"
		} else if question.Group != nil && question.Group.Type != "" {
			questionType = question.Group.Type
		}

		questionOut := QuestionOut{
			Content: question.Content[languageCode],
			Type:    questionType,
		}
		questionsOut = append(questionsOut, questionOut)
	}

	return questionsOut
}

func (env *QuestionAPI) AddTranslation(_ *gin.Context, questionInput *AddTranslationInput) error {
	var (
		questionID = questionInput.ID
		question   = questionInput.QuestionTranslationIn
		gameName   = questionInput.Name
		lang       = questionInput.Language
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

func (env *QuestionAPI) EnableQuestion(_ *gin.Context, questionInput *QuestionInput) (struct{}, error) {
	return env.updateEnable(questionInput, true)
}

func (env *QuestionAPI) DisableQuestion(_ *gin.Context, questionInput *QuestionInput) (struct{}, error) {
	return env.updateEnable(questionInput, false)
}

func (env *QuestionAPI) GetAllGroups(_ *gin.Context, groupInput *GroupInput) ([]string, error) {
	gameName := groupInput.Name
	round := groupInput.Round
	questionLogger := env.Logger.WithFields(log.Fields{
		"game_name": gameName,
		"round":     round,
	})

	questionLogger.Debug("Trying to get all groups")
	q := QuestionService{
		DB:       env.DB,
		GameName: gameName,
	}

	groups, err := q.GetGroups(round)
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
		questionID = questionInput.ID
		gameName   = questionInput.Name
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
