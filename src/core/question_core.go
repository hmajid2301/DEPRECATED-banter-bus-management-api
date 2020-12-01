package core

import (
	"fmt"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/games"
	"banter-bus-server/src/core/models"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// AddQuestion is add questions to a game.
func AddQuestion(gameName string, question models.GenericQuestion) error {
	gameType, err := validateAndGetGameType(gameName, question)
	if err != nil {
		return err
	}

	questionPath := gameType.GetQuestionPath()
	err = validateQuestionNotFound(gameName, questionPath, question.LanguageCode, question.Content)
	if err != nil {
		return err
	}

	t := true
	questionToAdd := map[string]models.Question{
		questionPath: {
			Content: map[string]string{
				question.LanguageCode: question.Content,
			},
			Enabled: &t,
		},
	}

	filter := &models.Game{Name: gameName}
	updated, err := database.AppendToEntry("game", filter, questionToAdd)
	if !updated || err != nil {
		errorMessage := "Failed to add a new question."
		log.Error(errorMessage)
		return errors.Errorf(errorMessage)
	}

	return nil
}

// UpdateQuestion is add questions to a game.
func UpdateQuestion(
	gameName string,
	existingQuestion models.GenericQuestion,
	questionContent string,
	questionLanguageCode string,
) error {
	gameType, err := validateAndGetGameType(gameName, existingQuestion)
	if err != nil {
		return err
	}

	questionPath := gameType.GetQuestionPath()

	originalQuestionExistsErr := validateQuestionFound(
		gameName,
		questionPath,
		existingQuestion.LanguageCode,
		existingQuestion.Content,
	)
	if originalQuestionExistsErr != nil {
		return originalQuestionExistsErr
	}

	newQuestionExistsErr := validateQuestionNotFound(gameName, questionPath, questionLanguageCode, questionContent)
	if newQuestionExistsErr != nil {
		return newQuestionExistsErr
	}

	filter := newQuestionFilter(questionPath, gameName, existingQuestion.Content, existingQuestion.LanguageCode)
	languagePath := fmt.Sprintf("content.%s", questionLanguageCode)
	questionToUpdate := newQuestion(questionPath, languagePath, questionContent)

	updated, err := database.UpsertEntry("game", filter, questionToUpdate)
	if !updated || err != nil {
		return errors.Errorf("Failed to update existing question.")
	}

	return nil
}

// RemoveQuestion is tp remove questions from a game.
func RemoveQuestion(gameName string, question models.GenericQuestion) error {
	gameType, err := validateAndGetGameType(gameName, question)
	if err != nil {
		return err
	}

	questionPath := gameType.GetQuestionPath()
	err = validateQuestionFound(gameName, questionPath, question.LanguageCode, question.Content)
	if err != nil {
		return err
	}

	questionToRemove := newEmtpyQuestion(questionPath, question.LanguageCode)
	filter := newQuestionFilter(questionPath, gameName, question.Content, question.LanguageCode)

	updated, err := database.RemoveEntry("game", filter, questionToRemove)
	if !updated || err != nil {
		return errors.Errorf("Failed to remove question.")
	}

	return nil
}

// UpdateEnableQuestion is used to update the enable state of a question.
func UpdateEnableQuestion(gameName string, enabled bool, question models.GenericQuestion) (bool, error) {
	gameType, err := validateAndGetGameType(gameName, question)
	if err != nil {
		return false, err
	}

	questionPath := gameType.GetQuestionPath()
	err = validateQuestionFound(gameName, questionPath, question.LanguageCode, question.Content)
	if err != nil {
		return false, err
	}

	filter := newQuestionFilter(questionPath, gameName, question.Content, "")
	update := newQuestion(questionPath, "enabled", enabled)

	updated, err := database.UpsertEntry("game", filter, update)
	if err != nil {
		return false, errors.Errorf("Failed to update question.")
	}
	return updated, err
}

func validateAndGetGameType(gameName string, question models.GenericQuestion) (games.PlayableGame, error) {
	game, _ := GetGame(gameName)
	if game.Name == "" {
		return nil, errors.NotFoundf("The game %s", gameName)
	}

	_, err := language.Parse(question.LanguageCode)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse language code: %s", question.LanguageCode), err)
		return nil, errors.BadRequestf("Invalid language code: %s", question.LanguageCode)
	}

	gameType, err := getGameType(gameName, question)
	if err != nil {
		return nil, err
	}

	err = gameType.ValidateQuestionInput()
	if err != nil {
		return nil, err
	}

	return gameType, nil
}

func validateQuestionNotFound(gameName string, questionPath string, languageCode string, content string) error {
	questionExists := doesQuestionExist(gameName, questionPath, languageCode, content)
	if questionExists {
		return errors.AlreadyExistsf("The question for game %s", gameName)
	}

	return nil
}

func validateQuestionFound(gameName string, questionPath string, languageCode string, content string) error {
	questionExists := doesQuestionExist(gameName, questionPath, languageCode, content)
	if !questionExists {
		return errors.NotFoundf("The question for game %s", gameName)
	}

	return nil
}

func doesQuestionExist(gameName string, questionPath string, languageCode string, content string) bool {
	var q *models.Game

	contentQuestionFilter := fmt.Sprintf("%s.content.%s", questionPath, languageCode)
	questionFilter := map[string]string{
		"name":                gameName,
		contentQuestionFilter: content,
	}

	err := database.Get("game", questionFilter, &q)
	return err == nil
}

func newQuestionFilter(
	questionPath string,
	gameName string,
	content string,
	optionalPath string,
) map[string]string {
	contentQuestionPath := fmt.Sprintf("%s.content", questionPath)
	if optionalPath != "" {
		contentQuestionPath += fmt.Sprintf(".%s", optionalPath)
	}
	filter := map[string]string{"name": gameName, contentQuestionPath: content}
	return filter
}

func newEmtpyQuestion(questionPath string, languageCode string) map[string]interface{} {
	languagePath := fmt.Sprintf("content.%s", languageCode)
	emptyObject := map[string]string{}
	emptyQuestion := newQuestion(questionPath, languagePath, emptyObject)
	return emptyQuestion
}

func newQuestion(questionPath string, attributePath string, content interface{}) map[string]interface{} {
	updatePath := fmt.Sprintf("%s.$.%s", questionPath, attributePath)
	update := map[string]interface{}{
		updatePath: content,
	}
	return update
}
