package games

import (
	"errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// Drawlosseum is the concrete type for the the game interface.
type Drawlosseum struct{}

// GetInfo is used to add new games of type `drawlosseum`.
func (d Drawlosseum) GetInfo(rulesURL string) models.GameInfo {
	t := true
	return models.GameInfo{
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

// GetQuestionPool gets the question pool structure for the Drawlosseum game.
func (d Drawlosseum) GetQuestionPool() interface{} {
	return models.DrawlosseumQuestionsPool{
		Drawings: []string{},
	}
}

// QuestionPoolToGenericQuestions converts question pool questions into generic questions that can be returned back to
// a client.
func (d Drawlosseum) QuestionPoolToGenericQuestions(questions interface{}) ([]models.GenericQuestion, error) {
	var newGenericQuestions []models.GenericQuestion
	drawlosseumQuestions, ok := questions.(models.DrawlosseumQuestionsPool)

	if !ok {
		errorMessage := "failed to convert type to DrawlosseumQuestionsPool"
		return []models.GenericQuestion{}, errors.New(errorMessage)
	}

	for _, question := range drawlosseumQuestions.Drawings {
		question := models.GenericQuestion{
			Content: question,
		}
		newGenericQuestions = append(newGenericQuestions, question)
	}

	return newGenericQuestions, nil
}
