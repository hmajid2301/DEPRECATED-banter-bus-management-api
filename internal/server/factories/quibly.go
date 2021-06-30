package factories

import (
	"github.com/juju/errors"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// Quibly struct which is the concrete type for game interface.
type Quibly struct{}

// ValidateQuestion is used to validate input for interacting with questions.
func (q Quibly) ValidateQuestion(question models.GenericQuestion) error {
	validRounds := map[string]bool{"pair": true, "group": true, "answers": true}
	if !validRounds[question.Round] {
		return errors.BadRequestf("invalid round %s", question.Round)
	}
	return nil
}

// NewServerStory returns "Quibly" story answers.
func (q Quibly) NewServerStory(story models.Story) (serverModels.Story, error) {
	storyAnswers, ok := story.Answers.(*models.StoryQuiblyAnswers)

	if !ok {
		return serverModels.Story{}, errors.Errorf("invalid answer for Quibly")
	}
	answers := newServerAnswersQuibly(storyAnswers)
	newStory := serverModels.Story{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswers: serverModels.StoryAnswers{
			Quibly: answers,
		},
	}
	return newStory, nil
}

func newServerAnswersQuibly(storyAnswers *models.StoryQuiblyAnswers) []serverModels.StoryQuibly {
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

// NewStory returns "Quibly" story answers.
func (q Quibly) NewStory(story serverModels.Story) (models.Story, error) {
	answers, err := newAnswersQuibly(story.StoryAnswers.Quibly)
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

func newAnswersQuibly(storyAnswers []serverModels.StoryQuibly) (models.StoryQuiblyAnswers, error) {
	var answers models.StoryQuiblyAnswers
	if len(storyAnswers) == 0 {
		return []models.StoryQuibly{}, errors.BadRequestf("no answers in the story.")
	}

	for _, storyAnswer := range storyAnswers {
		answer := models.StoryQuibly{
			Nickname: storyAnswer.Nickname,
			Answer:   storyAnswer.Answer,
			Votes:    storyAnswer.Votes,
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
