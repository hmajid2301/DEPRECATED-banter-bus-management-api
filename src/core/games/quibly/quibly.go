package quibly

import (
	"fmt"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/models"

	"github.com/juju/errors"
)

// Quibly type that implements PlayableGame.
type Quibly struct {
	CurrentQuestion models.GenericQuestion
}

// AddGame is used to add new games of type `quibly`.
func (q Quibly) AddGame(rulesURL string) (bool, error) {
	t := true
	var newGame = &models.Game{
		Name:     "quibly",
		RulesURL: rulesURL,
		Enabled:  &t,
		Questions: &models.QuiblyQuestions{
			Pair:   []models.Question{},
			Answer: []models.Question{},
			Group:  []models.Question{},
		},
	}

	inserted, err := database.Insert("game", newGame)
	return inserted, err
}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.pair".
func (q Quibly) GetQuestionPath() string {
	question := q.CurrentQuestion
	questionPath := fmt.Sprintf("questions.%s", question.Round)
	return questionPath
}

// ValidateQuestionInput is used to validate input for interacting with questions.
func (q Quibly) ValidateQuestionInput() error {
	question := q.CurrentQuestion
	validRounds := map[string]bool{"pair": true, "group": true, "answers": true}
	if !validRounds[question.Round] {
		return errors.BadRequestf("Invalid round %s", question.Round)
	}
	return nil
}
