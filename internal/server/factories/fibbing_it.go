package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// FibbingIt struct which is the concrete type for game interface.
type FibbingIt struct{}

// NewStory returns "FibbingIt" story answers.
func (f FibbingIt) NewStory(userStory models.Story) (serverModels.Story, error) {
	storyAnswer, ok := userStory.Answers.(*models.StoryFibbingItAnswers)
	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer for Fibbing It")
	}

	answers := newAnswersFibbingIt(storyAnswer)
	story := serverModels.Story{
		Question: userStory.Question,
		Round:    userStory.Round,
		StoryAnswers: serverModels.StoryAnswers{
			FibbingIt: answers,
		},
	}
	return story, nil
}

func newAnswersFibbingIt(storyAnswers *models.StoryFibbingItAnswers) []serverModels.StoryFibbingIt {
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
