package games

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
)

// Drawlosseum type that implements PlayableGame.
type Drawlosseum struct {
	DB core.Repository
}

// AddGame is used to add new games of type `drawlosseum`.
func (d Drawlosseum) AddGame(rulesURL string) (bool, error) {
	t := true
	var newGame = &models.GameInfo{
		Name:      "drawlosseum",
		RulesURL:  rulesURL,
		Enabled:   &t,
		Questions: &models.DrawlosseumQuestions{},
	}

	inserted, err := d.DB.Insert("game", newGame)
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
