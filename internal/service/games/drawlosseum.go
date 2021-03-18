package games

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Drawlosseum is the concrete type for the the game interface.
type Drawlosseum struct{}

// ValidateQuestion does nothing for the Drawlosseum game, as it requires no extra validation.
func (d Drawlosseum) ValidateQuestion(_ models.GenericQuestion) error {
	return nil
}
