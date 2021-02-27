package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Quibly struct which is the concrete type for game interface.
type Quibly struct{}

// NewQuestionPool returns "Quibly" pair, group and answers.
func (q Quibly) NewQuestionPool(questions models.QuestionPoolType) (serverModels.QuestionPoolQuestions, error) {
	quibly, ok := questions.(*models.QuiblyQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, errors.Errorf("invalid question for Quibly")
	}

	var pool serverModels.QuiblyQuestionsPool
	pool.Pair = quibly.Pair
	pool.Group = quibly.Group
	pool.Answers = quibly.Answers
	return serverModels.QuestionPoolQuestions{Quibly: pool}, nil
}

// NewStory returns "Quibly" story answers.
func (q Quibly) NewStory(userStory models.Story) (serverModels.Story, error) {
	storyAnswers, ok := userStory.Answers.(*models.StoryQuiblyAnswers)

	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer for Quibly")
	}
	answers := newAnswersQuibly(storyAnswers)
	story := serverModels.Story{
		Question: userStory.Question,
		Round:    userStory.Round,
		StoryAnswers: serverModels.StoryAnswers{
			Quibly: answers,
		},
	}
	return story, nil
}

func newAnswersQuibly(storyAnswers *models.StoryQuiblyAnswers) []serverModels.StoryQuibly {
	var answers []serverModels.StoryQuibly
	for _, storyAnswer := range *storyAnswers {
		answer := serverModels.StoryQuibly{
			Nickname: storyAnswer.Nickname,
			Answer:   storyAnswer.Answer,
			Votes:    storyAnswer.Votes,
		}
		answers = append(answers, answer)
	}

	return answers
}
