package service

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/juju/errors"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/language"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// QuestionService is struct data required by all question service functions.
type QuestionService struct {
	DB       database.Database
	GameName string
	Question models.GenericQuestion
}

// Add is add questions to a game.
func (q *QuestionService) Add() error {
	err := q.validateQuestion()
	if err != nil {
		return err
	}

	path, err := q.getQuestionPath()
	if err != nil {
		return err
	}

	err = q.validateQuestionNotFound(path, q.Question.LanguageCode, q.Question.Content)
	if err != nil {
		return err
	}

	question := newQuestion(q.Question, path)
	filter := map[string]string{"name": q.GameName}
	updated, err := question.AddToList(q.DB, filter)

	if !updated || err != nil {
		return errors.Errorf("failed to add a new question")
	}

	return nil
}

// AddTranslation is used to add a new question in a different language to a game.
func (q *QuestionService) AddTranslation(content string, langCode string) error {
	err := q.validateQuestion()
	if err != nil {
		return err
	}

	path, err := q.getQuestionPath()
	if err != nil {
		return err
	}

	currLangCode, currContent := q.Question.LanguageCode, q.Question.Content
	err = q.validateQuestionFound(path, currLangCode, currContent)
	if err != nil {
		return err
	}

	err = q.validateQuestionNotFound(path, langCode, content)
	if err != nil {
		return err
	}

	filter := q.getFilter(path, currContent, currLangCode)
	endPath := fmt.Sprintf("content.%s", langCode)
	question := getUpdatedQuestion(path, endPath, content)

	updated, err := q.DB.UpdateObject("game", filter, question)
	if !updated || err != nil {
		return errors.Errorf("failed to add translation to question %s", currContent)
	}

	return nil
}

// Remove removes questions from a game.
func (q *QuestionService) Remove() error {
	err := q.validateQuestion()
	if err != nil {
		return err
	}

	path, err := q.getQuestionPath()
	if err != nil {
		return err
	}
	err = q.validateQuestionFound(path, q.Question.LanguageCode, q.Question.Content)
	if err != nil {
		return err
	}

	question := emptyQuestion(path, q.Question.LanguageCode)
	filter := q.getFilter(path, q.Question.Content, q.Question.LanguageCode)

	updated, err := q.DB.RemoveObject("game", filter, question)
	if !updated || err != nil {
		return errors.Errorf("failed to remove question")
	}

	return nil
}

// UpdateEnable is used to update the enable state of a question.
func (q *QuestionService) UpdateEnable(
	enabled bool,
) (bool, error) {
	err := q.validateQuestion()
	if err != nil {
		return false, err
	}

	path, err := q.getQuestionPath()
	if err != nil {
		return false, err
	}

	err = q.validateQuestionFound(path, q.Question.LanguageCode, q.Question.Content)
	if err != nil {
		return false, err
	}

	filter := q.getFilter(path, q.Question.Content, "")
	question := getUpdatedQuestion(path, "enabled", enabled)

	updated, err := q.DB.UpdateObject("game", filter, question)
	if err != nil {
		return false, errors.Errorf("failed to update question")
	}
	return updated, err
}

// GetGroups gets all question group names for a given game and round.
func (q *QuestionService) GetGroups(round string) ([]string, error) {
	g := GameService{DB: q.DB, Name: q.GameName}
	game, err := g.Get()
	if err != nil {
		return nil, err
	}

	if !game.HasGroups(round) {
		return nil, errors.NotFoundf("cannot get question groups from round %s of game %s", round, q.GameName)
	}

	bytes, err := bson.MarshalExtJSON(game.Questions, true, true)
	if err != nil {
		return nil, err
	}

	var questions map[string]interface{}
	err = json.Unmarshal(bytes, &questions)
	if err != nil {
		return nil, err
	}

	groupInterface, roundPresent := questions[round]
	if !roundPresent || questions[round] == nil {
		return nil, errors.NotFoundf("cannot find round: %s", round)
	}
	groupList, err := getGroupList(groupInterface)
	return groupList, err
}

func (q *QuestionService) validateQuestion() error {
	gameService := GameService{DB: q.DB, Name: q.GameName}
	gameInfo, err := gameService.Get()

	if gameInfo.Name == "" {
		return errors.NotFoundf("game %s", q.GameName)
	} else if err != nil {
		return err
	}

	_, err = language.Parse(q.Question.LanguageCode)
	if err != nil {
		return errors.BadRequestf("invalid language code %s", q.Question.LanguageCode)
	}

	game, err := games.GetGame(q.GameName)
	if err != nil {
		return err
	}

	err = game.ValidateQuestion(q.Question)
	if err != nil {
		return err
	}

	return nil
}

func (q *QuestionService) getQuestionPath() (string, error) {
	game, err := games.GetGame(q.GameName)
	if err != nil {
		return "", err
	}

	path := game.GetQuestionPath(q.Question)
	return path, nil
}

func (q *QuestionService) validateQuestionNotFound(
	path string,
	langCode string,
	content string,
) error {
	questionExists := q.questionExist(path, langCode, content)
	if questionExists {
		return errors.AlreadyExistsf("The question for game %s", q.GameName)
	}

	return nil
}

func (q *QuestionService) validateQuestionFound(
	path string,
	langCode string,
	content string,
) error {
	questionExists := q.questionExist(path, langCode, content)
	if !questionExists {
		return errors.NotFoundf("the question for game %s", q.GameName)
	}

	return nil
}

func (q *QuestionService) questionExist(
	path string,
	langCode string,
	content string,
) bool {
	game := &models.Game{}

	contentFilter := fmt.Sprintf("%s.content.%s", path, langCode)
	questionFilter := map[string]string{
		"name":        q.GameName,
		contentFilter: content,
	}

	err := game.Get(q.DB, questionFilter)
	return err == nil
}

func newQuestion(question models.GenericQuestion, questionPath string) models.NewQuestion {
	t := true
	newQuestion := models.NewQuestion{
		questionPath: {
			Content: map[string]string{
				question.LanguageCode: question.Content,
			},
			Enabled: &t,
		},
	}

	return newQuestion
}

func (q *QuestionService) getFilter(
	path string,
	content string,
	optPath string,
) map[string]string {
	questionPath := fmt.Sprintf("%s.content", path)
	if optPath != "" {
		questionPath += fmt.Sprintf(".%s", optPath)
	}
	filter := map[string]string{"name": q.GameName, questionPath: content}
	return filter
}

func emptyQuestion(path string, langCode string) map[string]interface{} {
	languagePath := fmt.Sprintf("content.%s", langCode)
	empty := map[string]string{}
	emptyQuestion := getUpdatedQuestion(path, languagePath, empty)
	return emptyQuestion
}

func getUpdatedQuestion(path string, attrPath string, content interface{}) map[string]interface{} {
	questionPath := fmt.Sprintf("%s.$.%s", path, attrPath)
	question := map[string]interface{}{
		questionPath: content,
	}
	return question
}

func getGroupList(group interface{}) ([]string, error) {
	groups, ok := group.(map[string]interface{})
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
