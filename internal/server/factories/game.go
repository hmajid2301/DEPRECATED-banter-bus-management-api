package factories

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// Game is the interface for all games.
type Game interface {
	NewQuestionPool(questions interface{}) (serverModels.QuestionPoolQuestions, error)
	NewStory(story models.Story) serverModels.Story
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
