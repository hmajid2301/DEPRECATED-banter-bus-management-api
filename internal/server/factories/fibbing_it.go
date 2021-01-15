package factories

import (
	"fmt"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// FibbingIt struct which is the concrete type for game interface.
type FibbingIt struct{}

// NewQuestionPool returns "FibbingIt" free form, likely and opinion.
func (f FibbingIt) NewQuestionPool(questions interface{}) (serverModels.QuestionPoolQuestions, error) {
	fibbingItQuestions, ok := questions.(models.FibbingItQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, fmt.Errorf("unable to cast to FibbingItQuestionsPool")
	}

	var newFibbingItQuestionPool serverModels.FibbingItQuestionsPool
	newFibbingItQuestionPool.Opinion = fibbingItQuestions.Opinion
	newFibbingItQuestionPool.Likely = fibbingItQuestions.Likely
	newFibbingItQuestionPool.FreeForm = fibbingItQuestions.FreeForm
	return serverModels.QuestionPoolQuestions{FibbingIt: newFibbingItQuestionPool}, nil
}

// NewStory returns "FibbingIt" story answers.
func (f FibbingIt) NewStory(story models.Story) serverModels.Story {
	fibbingItAnswers := newAnswersFibbingIt(story.Answers)
	newFibbingItStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			FibbingIt: fibbingItAnswers,
		},
	}
	return newFibbingItStory
}

func newAnswersFibbingIt(answers interface{}) []serverModels.StoryFibbingIt {
	fibbingItAnswers := answers.([]models.StoryFibbingIt)

	var newAnswersFibbingIt []serverModels.StoryFibbingIt
	for _, answer := range fibbingItAnswers {
		newAnswer := serverModels.StoryFibbingIt{
			Nickname: answer.Nickname,
			Answer:   answer.Answer,
		}
		newAnswersFibbingIt = append(newAnswersFibbingIt, newAnswer)
	}

	return newAnswersFibbingIt
}
