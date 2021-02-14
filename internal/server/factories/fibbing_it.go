package factories

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// FibbingIt struct which is the concrete type for game interface.
type FibbingIt struct{}

// NewQuestionPool returns "FibbingIt" free form, likely and opinion.
func (f FibbingIt) NewQuestionPool(questions models.QuestionPoolType) (serverModels.QuestionPoolQuestions, error) {
	fibbingIt, ok := questions.(*models.FibbingItQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, errors.Errorf("invalid question type for game fibbing_it")
	}
	var newFibbingItQuestionPool serverModels.FibbingItQuestionsPool
	newFibbingItQuestionPool.Opinion = fibbingIt.Opinion
	newFibbingItQuestionPool.Likely = fibbingIt.Likely
	newFibbingItQuestionPool.FreeForm = fibbingIt.FreeForm
	return serverModels.QuestionPoolQuestions{FibbingIt: newFibbingItQuestionPool}, nil
}

// NewStory returns "FibbingIt" story answers.
func (f FibbingIt) NewStory(story models.Story) (serverModels.Story, error) {
	answers, ok := story.Answers.(*models.StoryFibbingItAnswers)
	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer type for game fibbing_it")
	}

	fibbingItAnswers := newAnswersFibbingIt(answers)
	newFibbingItStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			FibbingIt: fibbingItAnswers,
		},
	}
	return newFibbingItStory, nil
}

func newAnswersFibbingIt(answers *models.StoryFibbingItAnswers) []serverModels.StoryFibbingIt {
	var newAnswersFibbingIt []serverModels.StoryFibbingIt
	for _, answer := range *answers {
		newAnswer := serverModels.StoryFibbingIt{
			Nickname: answer.Nickname,
			Answer:   answer.Answer,
		}
		newAnswersFibbingIt = append(newAnswersFibbingIt, newAnswer)
	}

	return newAnswersFibbingIt
}
