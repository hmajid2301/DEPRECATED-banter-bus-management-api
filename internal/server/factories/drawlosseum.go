package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Drawlosseum struct which is the concrete type for game interface.
type Drawlosseum struct{}

// NewStory returns "Drawlosseum" style answers.
func (d Drawlosseum) NewStory(story models.Story) (serverModels.Story, error) {
	storyAnswers, ok := story.Answers.(*models.StoryDrawlosseumAnswers)
	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer for Drawlosseum")
	}
	answers := newDrawlosseumAnswers(storyAnswers)
	newStory := serverModels.Story{
		Question: story.Question,
		Nickname: story.Nickname,
		StoryAnswers: serverModels.StoryAnswers{
			Drawlosseum: answers,
		},
	}
	return newStory, nil
}

func newDrawlosseumAnswers(storyAnswers *models.StoryDrawlosseumAnswers) []serverModels.StoryDrawlosseum {
	var answers []serverModels.StoryDrawlosseum
	for _, storyAnswer := range *storyAnswers {
		answer := serverModels.StoryDrawlosseum{
			Color: storyAnswer.Color,
			Start: serverModels.DrawlosseumDrawingPoint{
				X: storyAnswer.Start.X,
				Y: storyAnswer.Start.Y,
			},
			End: serverModels.DrawlosseumDrawingPoint{
				X: storyAnswer.End.X,
				Y: storyAnswer.End.Y,
			},
		}
		answers = append(answers, answer)
	}

	return answers
}
