package games

import (
	"fmt"

	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// FibbingIt is the concrete type for the the game interface.
type FibbingIt struct{}

// GetInfo is used to add new games of type `fibbing_it`.
func (f FibbingIt) GetInfo(rulesURL string) models.GameInfo {
	t := true
	return models.GameInfo{
		Name:     "fibbing_it",
		RulesURL: rulesURL,
		Enabled:  &t,
		Questions: &models.FibbingItQuestions{
			Opinion:  map[string]map[string][]models.Question{},
			FreeForm: map[string][]models.Question{},
			Likely:   []models.Question{},
		},
	}
}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.likely".
func (f FibbingIt) GetQuestionPath(question models.GenericQuestion) string {
	questionPath := fmt.Sprintf("questions.%s", question.Round)

	if question.Group.Name != "" {
		questionPath += fmt.Sprintf(".%s", question.Group.Name)
	}
	if question.Group.Type != "" {
		questionPath += fmt.Sprintf(".%s", question.Group.Type)
	}
	return questionPath
}

// ValidateQuestionInput is used to validate input for interacting with questions.
func (f FibbingIt) ValidateQuestionInput(question models.GenericQuestion) error {
	validRounds := map[string]bool{"opinion": true, "likely": true, "free_form": true}
	validTypes := map[string]bool{"answers": true, "questions": true}

	if !validRounds[question.Round] {
		return errors.BadRequestf("Invalid round %s", question.Round)
	}

	if question.Group == nil {
		return errors.BadRequestf("Missing group information %s", question.Group)
	} else if question.Group.Type != "" && !validTypes[question.Group.Type] {
		return errors.BadRequestf("Invalid group type %s", question.Group.Type)
	}

	return nil
}
