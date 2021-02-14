package games

import (
	"fmt"

	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// Quibly type that implements PlayableGame.
type Quibly struct{}

// NewGame is gets the data for an empty/new `quibly` game.
func (q Quibly) NewGame(rulesURL string) models.Game {
	t := true
	return models.Game{
		Name:     "quibly",
		RulesURL: rulesURL,
		Enabled:  &t,
		Questions: &models.QuiblyQuestions{
			Pair:    []models.Question{},
			Answers: []models.Question{},
			Group:   []models.Question{},
		},
	}
}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.pair".
func (q Quibly) GetQuestionPath(question models.GenericQuestion) string {
	questionPath := fmt.Sprintf("questions.%s", question.Round)
	return questionPath
}

// ValidateQuestionInput is used to validate input for interacting with questions.
func (q Quibly) ValidateQuestionInput(question models.GenericQuestion) error {
	validRounds := map[string]bool{"pair": true, "group": true, "answers": true}
	if !validRounds[question.Round] {
		return errors.BadRequestf("Invalid round %s", question.Round)
	}
	return nil
}

// NewQuestionPool gets the question pool structure for the Quibly game.
func (q Quibly) NewQuestionPool() models.QuestionPoolType {
	quiblyQuestionPool := &models.QuiblyQuestionsPool{}
	quiblyQuestionPool.EmptyPoolQuestions()
	return quiblyQuestionPool
}

// QuestionPoolToGenericQuestions converts question pool questions into generic questions that can be returned back to
// a client.
func (q Quibly) QuestionPoolToGenericQuestions(questions models.QuestionPoolType) ([]models.GenericQuestion, error) {
	quibly, ok := questions.(*models.QuiblyQuestionsPool)
	if !ok {
		return nil, errors.Errorf("invalid question type for game quibly")
	}
	var newGenericQuestions []models.GenericQuestion

	questionsGroup := map[string]interface{}{
		"pair":    quibly.Pair,
		"answers": quibly.Answers,
		"group":   quibly.Group,
	}

	for round, group := range questionsGroup {
		for _, content := range group.([]string) {
			question := models.GenericQuestion{
				Content: content,
				Round:   round,
			}
			newGenericQuestions = append(newGenericQuestions, question)
		}
	}

	return newGenericQuestions, nil
}
