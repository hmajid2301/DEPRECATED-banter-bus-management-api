package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// FibbingIt struct which is the concrete type for game interface.
type FibbingIt struct{}

// NewServerStory returns "FibbingIt" story answers.
func (f FibbingIt) NewServerStory(story models.Story) (serverModels.Story, error) {
	storyAnswer, ok := story.Answers.(*models.StoryFibbingItAnswers)
	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer for Fibbing It")
	}

	answers := newServevAnswersFibbingIt(storyAnswer)
	newStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			FibbingIt: answers,
		},
	}
	return newStory, nil
}

func newServevAnswersFibbingIt(storyAnswers *models.StoryFibbingItAnswers) []serverModels.StoryFibbingIt {
	var answers []serverModels.StoryFibbingIt
	for _, storyAnswer := range *storyAnswers {
		answer := serverModels.StoryFibbingIt{
			Nickname: storyAnswer.Nickname,
			Answer:   storyAnswer.Answer,
		}
		answers = append(answers, answer)
	}

	return answers
}

// NewServerStory returns "FibbingIt" story answers.
func (f FibbingIt) NewStory(story serverModels.Story) (models.Story, error) {
	answers, err := newAnswersFibbingIt(story.StoryAnswers.FibbingIt)
	if err != nil {
		return models.Story{}, err
	}

	newStory := models.Story{
		Question: story.Question,
		Round:    story.Round,
		Answers:  answers,
	}
	return newStory, nil
}

func newAnswersFibbingIt(storyAnswers []serverModels.StoryFibbingIt) (models.StoryFibbingItAnswers, error) {
	var answers []models.StoryFibbingIt
	if len(storyAnswers) == 0 {
		return []models.StoryFibbingIt{}, errors.BadRequestf("no answers in the story.")
	}

	for _, storyAnswer := range storyAnswers {
		answer := models.StoryFibbingIt{
			Nickname: storyAnswer.Nickname,
			Answer:   storyAnswer.Answer,
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
