package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Drawlosseum struct which is the concrete type for game interface.
type Drawlosseum struct{}

// NewServerStory returns "Drawlosseum" style answers.
func (d Drawlosseum) NewServerStory(story models.Story) (serverModels.Story, error) {
	storyAnswers, ok := story.Answers.(*models.StoryDrawlosseumAnswers)
	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer for Drawlosseum")
	}
	answers := newServerDrawlosseumAnswers(storyAnswers)
	newStory := serverModels.Story{
		Question: story.Question,
		Nickname: story.Nickname,
		StoryAnswers: serverModels.StoryAnswers{
			Drawlosseum: answers,
		},
	}
	return newStory, nil
}

func newServerDrawlosseumAnswers(storyAnswers *models.StoryDrawlosseumAnswers) []serverModels.StoryDrawlosseum {
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

// NewStory returns "Drawlosseum" style answers.
func (d Drawlosseum) NewStory(story serverModels.Story) (models.Story, error) {
	answers, err := newDrawlosseumAnswers(story.StoryAnswers.Drawlosseum)
	if err != nil {
		return models.Story{}, err
	}

	newStory := models.Story{
		Question: story.Question,
		Nickname: story.Nickname,
		Answers:  answers,
	}
	return newStory, nil
}

func newDrawlosseumAnswers(storyAnswers []serverModels.StoryDrawlosseum) (models.StoryDrawlosseumAnswers, error) {
	var answers models.StoryDrawlosseumAnswers
	if len(storyAnswers) == 0 {
		return []models.StoryDrawlosseum{}, errors.BadRequestf("no answers in the story.")
	}

	for _, storyAnswer := range storyAnswers {
		answer := models.StoryDrawlosseum{
			Color: storyAnswer.Color,
			Start: models.DrawlosseumDrawingPoint{
				X: storyAnswer.Start.X,
				Y: storyAnswer.Start.Y,
			},
			End: models.DrawlosseumDrawingPoint{
				X: storyAnswer.End.X,
				Y: storyAnswer.End.Y,
			},
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
