package factories

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// Quibly struct which is the concrete type for game interface.
type Quibly struct{}

// NewQuestionPool returns "Quibly" pair, group and answers.
func (q Quibly) NewQuestionPool(questions models.QuestionPoolType) (serverModels.QuestionPoolQuestions, error) {
	quibly, ok := questions.(*models.QuiblyQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, errors.Errorf("invalid question type for game quibly")
	}

	var newQuiblyQuestionPool serverModels.QuiblyQuestionsPool
	newQuiblyQuestionPool.Pair = quibly.Pair
	newQuiblyQuestionPool.Group = quibly.Group
	newQuiblyQuestionPool.Answers = quibly.Answers
	return serverModels.QuestionPoolQuestions{Quibly: newQuiblyQuestionPool}, nil
}

// NewStory returns "Quibly" story answers.
func (q Quibly) NewStory(story models.Story) (serverModels.Story, error) {
	answers, ok := story.Answers.(*models.StoryQuiblyAnswers)

	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer type for game quibly")
	}
	quiblyAnswers := newAnswersQuibly(answers)
	newQuiblyStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			Quibly: quiblyAnswers,
		},
	}
	return newQuiblyStory, nil
}

func newAnswersQuibly(answers *models.StoryQuiblyAnswers) []serverModels.StoryQuibly {
	var newAnswersQuibly []serverModels.StoryQuibly
	for _, answer := range *answers {
		newAnswer := serverModels.StoryQuibly{
			Nickname: answer.Nickname,
			Answer:   answer.Answer,
			Votes:    answer.Votes,
		}
		newAnswersQuibly = append(newAnswersQuibly, newAnswer)
	}

	return newAnswersQuibly
}
