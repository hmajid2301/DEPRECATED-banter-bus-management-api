package factories

import (
	"fmt"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// Drawlosseum struct which is the concrete type for game interface.
type Drawlosseum struct{}

// NewQuestionPool returns "Drawlosseum" drawings.
func (d Drawlosseum) NewQuestionPool(questions interface{}) (serverModels.QuestionPoolQuestions, error) {
	drawlosseumQuestions, ok := questions.(models.DrawlosseumQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, fmt.Errorf("unable to cast to DrawlosseumQuestionsPool")
	}

	var newDrawlosseumQuestionPool serverModels.DrawlosseumQuestionsPool
	newDrawlosseumQuestionPool.Drawings = drawlosseumQuestions.Drawings
	return serverModels.QuestionPoolQuestions{Drawlosseum: newDrawlosseumQuestionPool}, nil
}

// NewStory returns "Drawlosseum" style answers.
func (d Drawlosseum) NewStory(story models.Story) serverModels.Story {
	drawlosseumAnswers := newDrawlosseumAnswers(story.Answers)
	newDrawlosseumStory := serverModels.Story{
		Question: story.Question,
		Nickname: story.Nickname,
		StoryAnswers: serverModels.StoryAnswers{
			Drawlosseum: drawlosseumAnswers,
		},
	}
	return newDrawlosseumStory
}

func newDrawlosseumAnswers(answers interface{}) []serverModels.StoryDrawlosseum {
	drawlosseumAnswers := answers.([]models.StoryDrawlosseum)

	var newAnswersDrawlosseum []serverModels.StoryDrawlosseum
	for _, answer := range drawlosseumAnswers {
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
