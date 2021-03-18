package games

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Quibly type that implements PlayableGame.
type Quibly struct{}

// ValidateQuestion is used to validate input for interacting with questions.
func (q Quibly) ValidateQuestion(question models.GenericQuestion) error {
	validRounds := map[string]bool{"pair": true, "group": true, "answers": true}
	if !validRounds[question.Round] {
		return errors.BadRequestf("invalid round %s", question.Round)
	}
	return nil
}
