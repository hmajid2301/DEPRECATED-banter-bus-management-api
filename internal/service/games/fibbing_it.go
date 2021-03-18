package games

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// FibbingIt is the concrete type for the the game interface.
type FibbingIt struct{}

// ValidateQuestion is used to validate input for interacting with questions.
func (f FibbingIt) ValidateQuestion(question models.GenericQuestion) error {
	validRounds := map[string]bool{"opinion": true, "likely": true, "free_form": true}
	validTypes := map[string]bool{"answer": true, "question": true}

	if !validRounds[question.Round] {
		return errors.BadRequestf("invalid round %s", question.Round)
	}

	if question.Group == nil {
		return errors.BadRequestf("missing group information %s", question.Group)
	} else if question.Group.Type != "" && !validTypes[question.Group.Type] {
		return errors.BadRequestf("invalid group type %s", question.Group.Type)
	}

	return nil
}
