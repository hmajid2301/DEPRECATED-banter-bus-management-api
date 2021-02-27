package games

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// GetGame is the factory function which will return the game struct based on the name.
func GetGame(name string) (models.Gamer, error) {
	switch name {
	case "quibly":
		return Quibly{}, nil
	case "fibbing_it":
		return FibbingIt{}, nil
	case "drawlosseum":
		return Drawlosseum{}, nil
	default:
		return nil, errors.NotFoundf("Game %s", name)
	}
}
