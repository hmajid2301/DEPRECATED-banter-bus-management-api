package games

import (
	"fmt"

	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// Quibly type that implements PlayableGame.
type Quibly struct{}

// GetInfo is used to add new games of type `quibly`.
func (q Quibly) GetInfo(rulesURL string) models.GameInfo {
	t := true
	return models.GameInfo{
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
