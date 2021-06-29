package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/juju/errors"
	"golang.org/x/text/language"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// QuestionService is struct data required by all question service functions.
type QuestionService struct {
	DB         database.Database
	GameName   string
	QuestionID string
	Question   models.GenericQuestion
}

// Add is add questions to a game.
func (q *QuestionService) Add() (string, error) {
	err := q.validateQuestion()
	if err != nil {
		return "", err
	}

	err = q.validateNotFound()
	if err != nil {
		return "", err
	}

	t := true
	uuidWithHyphen := uuid.New()
	uuid := strings.ReplaceAll(uuidWithHyphen.String(), "-", "")

	quest := models.Question{
		ID:       uuid,
		GameName: q.GameName,
		Round:    q.Question.Round,
		Enabled:  &t,
		Content: map[string]string{
			q.Question.LanguageCode: q.Question.Content,
		},
	}

	if q.Question.Group != nil {
		quest.Group.Name = q.Question.Group.Name
		quest.Group.Type = q.Question.Group.Type
	}

	inserted, err := quest.Add(q.DB)
	if !inserted || err != nil {
		return "", errors.Errorf("failed to add a new question %v", err)
	}

	return uuid, nil
}

// RemoveQuestion removes a question from a game.
func (q *QuestionService) RemoveQuestion() error {
	err := q.validateFound()
	if err != nil {
		return err
	}

	filter := q.filter()
	deleted, err := q.DB.Delete("question", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove question %v", err)
	}

	return nil
}

// AddTranslation is used to add a new question in a different language to a game.
func (q *QuestionService) AddTranslation(content string, langCode string) error {
	err := q.validateFound()
	if err != nil {
		return err
	}

	filter := q.filter()
	path := fmt.Sprintf("content.%s", langCode)
	translation := models.UpdateQuestion{
		path: content,
	}

	updated, err := translation.Add(q.DB, filter)
	if !updated || err != nil {
		return errors.Errorf("failed to add question translation %v", err)
	}

	return nil
}

// RemoveTranslation removes the content for one language code from a question.
func (q *QuestionService) RemoveTranslation(languageCode string) error {
	question, err := q.get()
	_, ok := question.Content[languageCode]
	if (err != nil) || !ok {
		return errors.NotFoundf("question with id %s and language code %s", q.QuestionID, languageCode)
	}

	filter := q.filter()
	path := fmt.Sprintf("content.%s", languageCode)
	translation := models.UpdateQuestion{
		path: "",
	}

	deleted, err := translation.Remove(q.DB, filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove question translation %v", err)
	}

	return nil
}

// UpdateEnable is used to update the enable state of a question.
func (q *QuestionService) UpdateEnable(enabled bool) (bool, error) {
	err := q.validateFound()
	if err != nil {
		return false, err
	}

	filter := q.filter()
	question := &models.Question{}
	err = question.Get(q.DB, filter)
	if err != nil {
		return false, err
	}

	question.Enabled = &enabled
	updated, err := question.Update(q.DB, filter)
	if err != nil {
		return false, errors.Errorf("failed to update question %v", err)
	}
	return updated, err
}

// GetGroups gets all question group names for a given game and round.
func (q *QuestionService) GetGroups(round string) ([]string, error) {
	gameService := GameService{DB: q.DB, Name: q.GameName}
	game, err := gameService.Get()
	if err != nil {
		return nil, err
	}

	if !game.HasGroups(round) {
		return nil, errors.NotFoundf("cannot get question groups from round %s of game %s", round, q.GameName)
	}

	filter := map[string]string{
		"game_name": q.GameName,
		"round":     round,
	}
	uniqGroups, err := q.DB.GetUnique("question", filter, "group.name")
	if err != nil {
		return nil, err
	}

	sort.Strings(uniqGroups)
	return uniqGroups, nil
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

func (q *QuestionService) validateNotFound() error {
	question, err := q.get()
	exists := (err == nil) || (question.Content != nil)
	if exists {
		return errors.AlreadyExistsf("the question '%s' for game %s", q.Question.Content, q.GameName)
	}

	return nil
}

func (q *QuestionService) validateFound() error {
	question, err := q.get()
	exists := (err == nil) || (question.Content != nil)
	if !exists {
		return errors.NotFoundf("the question with ID %s for game %s", q.QuestionID, q.GameName)
	}

	return nil
}

func (q *QuestionService) get() (*models.Question, error) {
	filter := q.filter()
	question := &models.Question{}
	err := question.Get(q.DB, filter)
	return question, err
}

func (q *QuestionService) filter() map[string]string {
	filter := map[string]string{
		"game_name": q.GameName,
	}

	if q.QuestionID != "" {
		filter["id"] = q.QuestionID
	}

	if q.Question.Content != "" {
		languageCode := "en"
		if q.Question.LanguageCode != "" {
			languageCode = q.Question.LanguageCode
		}

		questionPath := fmt.Sprintf("content.%s", languageCode)
		filter[questionPath] = q.Question.Content
	}
	return filter
}
