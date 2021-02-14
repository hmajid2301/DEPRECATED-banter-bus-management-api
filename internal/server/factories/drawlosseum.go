package factories

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// Drawlosseum struct which is the concrete type for game interface.
type Drawlosseum struct{}

// NewQuestionPool returns "Drawlosseum" drawings.
func (d Drawlosseum) NewQuestionPool(questions models.QuestionPoolType) (serverModels.QuestionPoolQuestions, error) {
	drawlosseum, ok := questions.(*models.DrawlosseumQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, errors.Errorf("invalid question type for game drawlosseum")
	}

	var newDrawlosseumQuestionPool serverModels.DrawlosseumQuestionsPool
	newDrawlosseumQuestionPool.Drawings = drawlosseum.Drawings
	return serverModels.QuestionPoolQuestions{Drawlosseum: newDrawlosseumQuestionPool}, nil
}

// NewStory returns "Drawlosseum" style answers.
func (d Drawlosseum) NewStory(story models.Story) (serverModels.Story, error) {
	answers, ok := story.Answers.(*models.StoryDrawlosseumAnswers)
	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer type for game drawlosseum")
	}
	drawlosseumAnswers := newDrawlosseumAnswers(answers)
	newDrawlosseumStory := serverModels.Story{
		Question: story.Question,
		Nickname: story.Nickname,
		StoryAnswers: serverModels.StoryAnswers{
			Drawlosseum: drawlosseumAnswers,
		},
	}
	return newDrawlosseumStory, nil
}

func newDrawlosseumAnswers(answers *models.StoryDrawlosseumAnswers) []serverModels.StoryDrawlosseum {
	var newAnswersDrawlosseum []serverModels.StoryDrawlosseum
	for _, answer := range *answers {
		newAnswer := serverModels.StoryDrawlosseum{
			Color: answer.Color,
			Start: serverModels.DrawlosseumDrawingPoint{
				X: answer.Start.X,
				Y: answer.Start.Y,
			},
			End: serverModels.DrawlosseumDrawingPoint{
				X: answer.End.X,
				Y: answer.End.Y,
			},
		}
		newAnswersDrawlosseum = append(newAnswersDrawlosseum, newAnswer)
	}

	return newAnswersDrawlosseum
}
