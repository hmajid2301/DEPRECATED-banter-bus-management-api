package factories

import (
	"fmt"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// Quibly struct which is the concrete type for game interface.
type Quibly struct{}

// NewQuestionPool returns "Quibly" pair, group and answers.
func (q Quibly) NewQuestionPool(questions interface{}) (serverModels.QuestionPoolQuestions, error) {
	quiblyQuestions, ok := questions.(models.QuiblyQuestionsPool)
	if !ok {
		return serverModels.QuestionPoolQuestions{}, fmt.Errorf("unable to cast to QuiblyQuestionsPool")
	}

	var newQuiblyQuestionPool serverModels.QuiblyQuestionsPool
	newQuiblyQuestionPool.Pair = quiblyQuestions.Pair
	newQuiblyQuestionPool.Group = quiblyQuestions.Group
	newQuiblyQuestionPool.Answers = quiblyQuestions.Answers
	return serverModels.QuestionPoolQuestions{Quibly: newQuiblyQuestionPool}, nil
}

// NewStory returns "Quibly" story answers.
func (q Quibly) NewStory(story models.Story) serverModels.Story {
	quiblyAnswers := newAnswersQuibly(story.Answers)
	newQuiblyStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			Quibly: quiblyAnswers,
		},
	}
	return newQuiblyStory
}

func newAnswersQuibly(answers interface{}) []serverModels.StoryQuibly {
	quiblyAnswers := answers.([]models.StoryQuibly)

	var newAnswersQuibly []serverModels.StoryQuibly
	for _, answer := range quiblyAnswers {
		newAnswer := serverModels.StoryQuibly{
			Nickname: answer.Nickname,
			Answer:   answer.Answer,
			Votes:    answer.Votes,
		}
		newAnswersQuibly = append(newAnswersQuibly, newAnswer)
	}

	return newAnswersQuibly
}
