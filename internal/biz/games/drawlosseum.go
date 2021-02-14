package games

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// Drawlosseum is the concrete type for the the game interface.
type Drawlosseum struct{}

// NewGame is gets the data for an empty/new `drawlosseum` game.
func (d Drawlosseum) NewGame(rulesURL string) models.Game {
	t := true
	return models.Game{
		Name:      "drawlosseum",
		RulesURL:  rulesURL,
		Enabled:   &t,
		Questions: &models.DrawlosseumQuestions{},
	}
}

// ValidateQuestionInput does nothing for the Drawlosseum game, as it requires no extra validation.
func (d Drawlosseum) ValidateQuestionInput(_ models.GenericQuestion) error {
	return nil
}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.drawings".
func (d Drawlosseum) GetQuestionPath(_ models.GenericQuestion) string {
	return "questions.drawings"
}

// NewQuestionPool gets the question pool structure for the Drawlosseum game.
func (d Drawlosseum) NewQuestionPool() models.QuestionPoolType {
	drawlosseumQuestionPool := &models.DrawlosseumQuestionsPool{}
	drawlosseumQuestionPool.EmptyPoolQuestions()
	return drawlosseumQuestionPool
}

// QuestionPoolToGenericQuestions converts question pool questions into generic questions that can be returned back to
// a client.
func (d Drawlosseum) QuestionPoolToGenericQuestions(
	questions models.QuestionPoolType,
) ([]models.GenericQuestion, error) {
	drawlosseum, ok := questions.(*models.DrawlosseumQuestionsPool)
	if !ok {
		return nil, errors.Errorf("invalid question type for game drawlosseum")
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
