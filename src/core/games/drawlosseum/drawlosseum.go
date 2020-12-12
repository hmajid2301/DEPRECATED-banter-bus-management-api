package drawlosseum

import (
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/models"
)

// Drawlosseum type that implements PlayableGame.
type Drawlosseum struct{}

// AddGame is used to add new games of type `drawlosseum`.
func (d Drawlosseum) AddGame(rulesURL string) (bool, error) {
	t := true
	var newGame = &models.Game{
		Name:      "drawlosseum",
		RulesURL:  rulesURL,
		Enabled:   &t,
		Questions: &models.DrawlosseumQuestions{},
	}

	inserted, err := database.Insert("game", newGame)
	return inserted, err
}

// ValidateQuestionInput does nothing for the Drawlosseum game, as it requires no extra validation.
func (d Drawlosseum) ValidateQuestionInput() error {
	return nil
}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.drawings".
func (d Drawlosseum) GetQuestionPath() string {
	return "questions.drawings"
}

// QuestionPoolToGenericQuestions converts question pool questions into generic questions that can be returned back to
// a client.
func (d Drawlosseum) QuestionPoolToGenericQuestions(questions interface{}) []models.GenericQuestion {
	var newGenericQuestions []models.GenericQuestion
	drawlosseumQuestions := questions.(models.DrawlosseumQuestionsPool)
	for _, question := range drawlosseumQuestions.Drawings {
		question := models.GenericQuestion{
			Content: question,
		}
		newGenericQuestions = append(newGenericQuestions, question)
	}

	return newGenericQuestions
}
