package biz

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/language"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// QuestionService is struct data required by all question service functions.
type QuestionService struct {
	DB models.Repository
}

// Add is add questions to a game.
func (q *QuestionService) Add(gameName string, question models.GenericQuestion) error {
	game, err := q.validateAndGetGame(gameName, question)
	if err != nil {
		return err
	}

	questionPath := game.GetQuestionPath(question)
	err = q.validateQuestionNotFound(gameName, questionPath, question.LanguageCode, question.Content)
	if err != nil {
		return err
	}

	t := true
	newQuestion := models.NewQuestion{
		questionPath: {
			Content: map[string]string{
				question.LanguageCode: question.Content,
			},
			Enabled: &t,
		},
	}

	filter := map[string]string{"name": gameName}
	updated, err := newQuestion.AddToList(q.DB, filter)
	if !updated || err != nil {
		return errors.Errorf("failed to add a new question")
	}

	return nil
}

// AddTranslation is used to add a new question in a different language to a game.
func (q *QuestionService) AddTranslation(
	gameName string,
	existingQuestion models.GenericQuestion,
	questionContent string,
	questionLanguageCode string,
) error {
	game, err := q.validateAndGetGame(gameName, existingQuestion)
	if err != nil {
		return err
	}

	questionPath := game.GetQuestionPath(existingQuestion)
	originalQuestionExistsErr := q.validateQuestionFound(
		gameName,
		questionPath,
		existingQuestion.LanguageCode,
		existingQuestion.Content,
	)
	if originalQuestionExistsErr != nil {
		return originalQuestionExistsErr
	}

	newQuestionExistsErr := q.validateQuestionNotFound(
		gameName,
		questionPath,
		questionLanguageCode,
		questionContent,
	)
	if newQuestionExistsErr != nil {
		return newQuestionExistsErr
	}

	filter := newQuestionFilter(questionPath, gameName, existingQuestion.Content, existingQuestion.LanguageCode)
	languagePath := fmt.Sprintf("content.%s", questionLanguageCode)
	questionToUpdate := newQuestion(questionPath, languagePath, questionContent)

	updated, err := q.DB.UpdateObject("game", filter, questionToUpdate)
	if !updated || err != nil {
		return errors.Errorf("failed to update existing question")
	}

	return nil
}

// Remove removes questions from a game.
func (q *QuestionService) Remove(gameName string, question models.GenericQuestion) error {
	game, err := q.validateAndGetGame(gameName, question)
	if err != nil {
		return err
	}

	questionPath := game.GetQuestionPath(question)
	err = q.validateQuestionFound(gameName, questionPath, question.LanguageCode, question.Content)
	if err != nil {
		return err
	}

	questionToRemove := newEmptyQuestion(questionPath, question.LanguageCode)
	filter := newQuestionFilter(questionPath, gameName, question.Content, question.LanguageCode)

	updated, err := q.DB.RemoveObject("game", filter, questionToRemove)
	if !updated || err != nil {
		return errors.Errorf("failed to remove question")
	}

	return nil
}

// UpdateEnable is used to update the enable state of a question.
func (q *QuestionService) UpdateEnable(
	gameName string,
	enabled bool,
	question models.GenericQuestion,
) (bool, error) {
	game, err := q.validateAndGetGame(gameName, question)
	if err != nil {
		return false, err
	}

	questionPath := game.GetQuestionPath(question)
	err = q.validateQuestionFound(gameName, questionPath, question.LanguageCode, question.Content)
	if err != nil {
		return false, err
	}

	filter := newQuestionFilter(questionPath, gameName, question.Content, "")
	update := newQuestion(questionPath, "enabled", enabled)

	updated, err := q.DB.UpdateObject("game", filter, update)
	if err != nil {
		return false, errors.Errorf("failed to update question")
	}
	return updated, err
}

// GetGroups gets all question group names for a given game and round.
func (q *QuestionService) GetGroups(gameName string, round string) ([]string, error) {
	gameService := GameService{DB: q.DB}
	game, err := gameService.Get(gameName)
	if err != nil {
		return nil, err
	}

	if !game.HasGroups(round) {
		return nil, errors.NotFoundf("cannot get question groups from round %s of game %s", round, gameName)
	}

	bytesData, err := bson.MarshalExtJSON(game.Questions, true, true)
	if err != nil {
		return nil, err
	}

	var questions map[string]interface{}

	err = json.Unmarshal(bytesData, &questions)
	if err != nil {
		return nil, err
	}

	groupInterface, roundPresent := questions[round]
	if !roundPresent || questions[round] == nil {
		return nil, errors.NotFoundf("cannot find round: %s", round)
	}

	groups, ok := groupInterface.(map[string]interface{})
	if !ok {
		return []string{}, errors.Errorf("failed to convert to type `map[string]interface{}`")
	}

	var groupList []string
	for group := range groups {
		groupList = append(groupList, group)
	}
	sort.Strings(groupList)
	return groupList, nil
}

func (q *QuestionService) validateAndGetGame(
	gameName string,
	question models.GenericQuestion,
) (models.IGame, error) {
	gameService := GameService{DB: q.DB}
	gameInfo, err := gameService.Get(gameName)
	if gameInfo.Name == "" {
		return nil, errors.NotFoundf("The game %s", gameName)
	} else if err != nil {
		return nil, err
	}

	_, err = language.Parse(question.LanguageCode)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse language code %s", question.LanguageCode), err)
		return nil, errors.BadRequestf("Invalid language code: %s", question.LanguageCode)
	}

	game, err := games.GetGame(gameName)
	if err != nil {
		return nil, err
	}

	err = game.ValidateQuestionInput(question)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (q *QuestionService) validateQuestionNotFound(
	gameName string,
	questionPath string,
	languageCode string,
	content string,
) error {
	questionExists := q.doesQuestionExist(gameName, questionPath, languageCode, content)
	if questionExists {
		return errors.AlreadyExistsf("The question for game %s", gameName)
	}

	return nil
}

func (q *QuestionService) validateQuestionFound(
	gameName string,
	questionPath string,
	languageCode string,
	content string,
) error {
	questionExists := q.doesQuestionExist(gameName, questionPath, languageCode, content)
	if !questionExists {
		return errors.NotFoundf("The question for game %s", gameName)
	}

	return nil
}

func (q *QuestionService) doesQuestionExist(
	gameName string,
	questionPath string,
	languageCode string,
	content string,
) bool {
	game := &models.Game{}

	contentQuestionFilter := fmt.Sprintf("%s.content.%s", questionPath, languageCode)
	questionFilter := map[string]string{
		"name":                gameName,
		contentQuestionFilter: content,
	}

	err := game.Get(q.DB, questionFilter)
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

func newEmptyQuestion(questionPath string, languageCode string) map[string]interface{} {
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
