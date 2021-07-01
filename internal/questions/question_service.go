package questions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

type QuestionService struct {
	DB         database.Database
	GameName   string
	QuestionID string
	Question   GenericQuestion
}

func (q *QuestionService) Add() (string, error) {
	err := q.validateNotFound()
	if err != nil {
		return "", err
	}

	t := true
	uuidWithHyphen := uuid.New()
	uuid := strings.ReplaceAll(uuidWithHyphen.String(), "-", "")

	question := Question{
		ID:       uuid,
		GameName: q.GameName,
		Round:    q.Question.Round,
		Enabled:  &t,
		Content: map[string]string{
			q.Question.LanguageCode: q.Question.Content,
		},
	}

	if q.Question.Group != nil {
		question.Group.Name = q.Question.Group.Name
		question.Group.Type = q.Question.Group.Type
	}

	inserted, err := question.Add(q.DB)
	if !inserted || err != nil {
		return "", errors.Errorf("failed to add a new question %v", err)
	}

	return uuid, nil
}

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

func (q *QuestionService) AddTranslation(content string, langCode string) error {
	err := q.validateFound()
	if err != nil {
		return err
	}

	filter := q.filter()
	path := fmt.Sprintf("content.%s", langCode)
	translation := UpdateQuestion{
		path: content,
	}

	updated, err := translation.Add(q.DB, filter)
	if !updated || err != nil {
		return errors.Errorf("failed to add question translation %v", err)
	}

	return nil
}

func (q *QuestionService) RemoveTranslation(languageCode string) error {
	question, err := q.get()
	_, ok := question.Content[languageCode]
	if (err != nil) || !ok {
		return errors.NotFoundf("question with id %s and language code %s", q.QuestionID, languageCode)
	}

	filter := q.filter()
	path := fmt.Sprintf("content.%s", languageCode)
	translation := UpdateQuestion{
		path: "",
	}

	deleted, err := translation.Remove(q.DB, filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove question translation %v", err)
	}

	return nil
}

func (q *QuestionService) UpdateEnable(enabled bool) (bool, error) {
	err := q.validateFound()
	if err != nil {
		return false, err
	}

	filter := q.filter()
	question := &Question{}
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

func (q *QuestionService) GetGroups(round string) ([]string, error) {
	game, err := GetGame(q.GameName)
	if err != nil {
		return []string{""}, err
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

func (q *QuestionService) get() (*Question, error) {
	filter := q.filter()
	question := &Question{}
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

type Questioner interface {
	ValidateQuestion(question QuestionIn) error
	HasGroups(_ string) bool
}

func GetGame(name string) (Questioner, error) {
	switch name {
	case "quibly":
		return Quibly{}, nil
	case "fibbing_it":
		return FibbingIt{}, nil
	case "drawlosseum":
		return Drawlosseum{}, nil
	default:
		return nil, errors.NotFoundf("Game %s", name)
	}
}
