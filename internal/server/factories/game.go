package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Game is the interface for all games.
type Game interface {
	NewQuestionPool(questions models.QuestionPoolType) (serverModels.QuestionPoolQuestions, error)
	NewStory(story models.Story) (serverModels.Story, error)
}

// GetGame is the factory function which returns the game struct depending on the name.
func GetGame(name string) (Game, error) {
	switch name {
	case "quibly":
		return Quibly{}, nil
	case "fibbing_it":
		return FibbingIt{}, nil
	case "drawlosseum":
		return Drawlosseum{}, nil
	}

	return nil, errors.NotFoundf("Game %s", name)
}
