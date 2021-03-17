package games

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Drawlosseum is the concrete type for the the game interface.
type Drawlosseum struct{}

// ValidateQuestion does nothing for the Drawlosseum game, as it requires no extra validation.
func (d Drawlosseum) ValidateQuestion(_ models.GenericQuestion) error {
	return nil
}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.drawings".
func (d Drawlosseum) GetQuestionPath(_ models.GenericQuestion) string {
	return "questions.drawings"
}

// NewQuestionPool gets the question pool structure for the Drawlosseum game.
func (d Drawlosseum) NewQuestionPool() models.QuestionPoolType {
	pool := &models.DrawlosseumQuestionsPool{}
	pool.NewPool()
	return pool
}

// QuestionPoolToGenericQuestions converts question pool questions into generic questions that can be returned back to
// a client.
func (d Drawlosseum) QuestionPoolToGenericQuestions(
	questions models.QuestionPoolType,
) ([]models.GenericQuestion, error) {
	drawlosseum, ok := questions.(*models.DrawlosseumQuestionsPool)
	if !ok {
		return nil, errors.Errorf("invalid question for Drawlosseum")
	}
	var newGenericQuestions []models.GenericQuestion

	for _, question := range drawlosseum.Drawings {
		question := models.GenericQuestion{
			Content: question,
		}
		newGenericQuestions = append(newGenericQuestions, question)
	}

	return newGenericQuestions, nil
}