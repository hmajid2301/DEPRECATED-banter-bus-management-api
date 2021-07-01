package story

import (
	"github.com/juju/errors"
)

type Quibly struct{}

func (q Quibly) NewStoryOut(story Story) (StoryInOut, error) {
	storyAnswers, ok := story.Answers.(*QuiblyAnswers)
	if !ok {
		return StoryInOut{}, errors.Errorf("invalid answer for Quibly")
	}

	answers := q.newAnswersOut(storyAnswers)
	newStory := StoryInOut{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswersInOut: StoryAnswersInOut{
			Quibly: answers,
		},
	}
	return newStory, nil
}

func (q Quibly) newAnswersOut(storyAnswers *QuiblyAnswers) QuiblyAnswersInOut {
	var answers QuiblyAnswersInOut
	for _, storyAnswer := range *storyAnswers {
		answer := QuiblyAnswerInOut(storyAnswer)
		answers = append(answers, answer)
	}

	return answers
}

func (q Quibly) NewStory(story StoryInOut) (Story, error) {
	answers, err := q.newAnswers(story.StoryAnswersInOut.Quibly)
	if err != nil {
		return Story{}, err
	}

	newStory := Story{
		Question: story.Question,
		Round:    story.Round,
		Answers:  answers,
	}
	return newStory, nil
}

func (q Quibly) newAnswers(storyAnswers QuiblyAnswersInOut) (QuiblyAnswers, error) {
	var answers QuiblyAnswers
	if len(storyAnswers) == 0 {
		return QuiblyAnswers{}, errors.BadRequestf("no answers in the story.")
	}

	for _, storyAnswer := range storyAnswers {
		answer := QuiblyAnswer(storyAnswer)
		answers = append(answers, answer)
	}

	return answers, nil
}
