package questions

import (
	"github.com/juju/errors"
)

type Quibly struct{}

func (q Quibly) ValidateQuestion(question QuestionIn) error {
	validRounds := map[string]bool{"pair": true, "group": true, "answers": true}
	round := question.Round
	if !validRounds[round] {
		return errors.BadRequestf("invalid round %s", round)
	}
	return nil
}

func (q Quibly) HasGroups(round string) bool {
	return false
}
