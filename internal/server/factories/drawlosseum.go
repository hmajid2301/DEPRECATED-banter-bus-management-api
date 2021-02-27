package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Drawlosseum struct which is the concrete type for game interface.
type Drawlosseum struct{}

// NewQuestionPool returns "Drawlosseum" drawings.
func (d Drawlosseum) NewQuestionPool(questions models.QuestionPoolType) (serverModels.QuestionPoolQuestions, error) {
	drawlosseum, ok := questions.(*models.DrawlosseumQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, errors.Errorf("invalid question for drawlosseum")
	}

	var pool serverModels.DrawlosseumQuestionsPool
	pool.Drawings = drawlosseum.Drawings
	return serverModels.QuestionPoolQuestions{Drawlosseum: pool}, nil
}

// NewStory returns "Drawlosseum" style answers.
func (d Drawlosseum) NewStory(userStory models.Story) (serverModels.Story, error) {
	storyAnswers, ok := userStory.Answers.(*models.StoryDrawlosseumAnswers)
	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer for Drawlosseum")
	}
	answers := newDrawlosseumAnswers(storyAnswers)
	story := serverModels.Story{
		Question: userStory.Question,
		Nickname: userStory.Nickname,
		StoryAnswers: serverModels.StoryAnswers{
			Drawlosseum: answers,
		},
	}
	return story, nil
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
