package service

import (
	"fmt"
	"sort"

	"github.com/juju/errors"
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

	err = q.validateNotFound()
	if err != nil {
		return err
	}

	t := true
	quest := models.Question{
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
		return errors.Errorf("failed to add a new question")
	}

	return nil
}

// RemoveQuestion removes a question from a game.
func (q *QuestionService) RemoveQuestion() error {
	err := q.validateQuestion()
	if err != nil {
		return err
	}

	err = q.validateFound()
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
	err := q.validateQuestion()
	if err != nil {
		return err
	}

	err = q.validateFound()
	if err != nil {
		return err
	}

	newQuestion := models.GenericQuestion{
		Content:      content,
		LanguageCode: langCode,
		Round:        q.Question.Round,
	}

	newQ := QuestionService{
		DB:       q.DB,
		GameName: q.GameName,
		Question: newQuestion,
	}
	err = newQ.validateNotFound()
	if err != nil {
		return err
	}

	filter := q.filter()
	path := fmt.Sprintf("content.%s", langCode)
	delete(filter, path)
	translation := models.UpdateQuestion{
		path: content,
	}

	updated, err := translation.Add(q.DB, filter)
	if !updated || err != nil {
		return errors.Errorf("failed to update existing question %v", err)
	}

	return nil
}

// RemoveTranslation removes the content for one language code from a question.
func (q *QuestionService) RemoveTranslation() error {
	err := q.validateQuestion()
	if err != nil {
		return err
	}

	err = q.validateFound()
	if err != nil {
		return err
	}

	filter := q.filter()
	path := fmt.Sprintf("content.%s", q.Question.LanguageCode)
	translation := models.UpdateQuestion{
		path: q.Question.Content,
	}

	deleted, err := translation.Remove(q.DB, filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove question")
	}

	return nil
}

// UpdateEnable is used to update the enable state of a question.
func (q *QuestionService) UpdateEnable(enabled bool) (bool, error) {
	err := q.validateQuestion()
	if err != nil {
		return false, err
	}

	filter := q.filter()
	currQuest := &models.Question{}
	err = currQuest.Get(q.DB, filter)
	if err != nil {
		return false, err
	}

	currQuest.Enabled = &enabled
	updated, err := currQuest.Update(q.DB, filter)
	if err != nil {
		return false, errors.Errorf("failed to update question")
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
	exists := q.exist()
	if exists {
		return errors.AlreadyExistsf("the question '%s' for game %s", q.Question.Content, q.GameName)
	}

	return nil
}

func (q *QuestionService) validateFound() error {
	exists := q.exist()
	if !exists {
		return errors.NotFoundf("the question '%s' for %s", q.Question.Content, q.GameName)
	}

	return nil
}

func (q *QuestionService) exist() bool {
	filter := q.filter()
	currQuest := &models.Question{}
	err := currQuest.Get(q.DB, filter)
	return (err == nil) || (currQuest.Content != nil)
}

func (q *QuestionService) filter() map[string]string {
	contentFilter := fmt.Sprintf("content.%s", q.Question.LanguageCode)
	filter := map[string]string{
		"game_name":   q.GameName,
		contentFilter: q.Question.Content,
	}

	if q.Question.Round != "" {
		filter["round"] = q.Question.Round
	}

	if q.Question.Group != nil {
		if q.Question.Group.Name != "" {
			filter["group.name"] = q.Question.Group.Name
		}

		if q.Question.Group.Type != "" {
			filter["group.type"] = q.Question.Group.Type
		}
	}
	return filter
}
