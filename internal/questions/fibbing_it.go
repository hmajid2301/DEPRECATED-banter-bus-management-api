package questions

import (
	"fmt"

	"github.com/juju/errors"
)

type FibbingIt struct{}

func (f FibbingIt) ValidateQuestion(question QuestionIn) error {
	validRounds := map[string]bool{"opinion": true, "likely": true, "free_form": true}
	validTypes := map[string]bool{"answer": true, "question": true}

	round := question.Round
	if !validRounds[round] {
		return errors.BadRequestf("invalid round %s", round)
	}

	group := question.Group
	if group != nil && round != "likely" {
		groupName := group.Name
		groupType := group.Type

		if groupName == "" {
			return errors.BadRequestf("missing group information")
		} else if groupType != "" && !validTypes[groupType] {
			return errors.BadRequestf("invalid group type %s", groupType)
		}
	}

	return nil
}

func (f FibbingIt) HasGroups(round string) bool {
	var gameRoundsWithGroups = map[string]struct{}{
		"fibbing_it.opinion":   {},
		"fibbing_it.free_form": {},
	}

	queryString := fmt.Sprintf("fibbing_it.%s", round)
	_, isPresent := gameRoundsWithGroups[queryString]
	return isPresent
}
