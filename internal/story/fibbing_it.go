package story

import (
	"github.com/juju/errors"
)

type FibbingIt struct{}

func (f FibbingIt) NewStoryOut(story Story) (StoryInOut, error) {
	storyAnswer, ok := story.Answers.(*FibbingItAnswers)
	if !ok {
		return StoryInOut{}, errors.Errorf("invalid answer for Fibbing It")
	}

	answers := f.newAnswersOut(storyAnswer)
	newStory := StoryInOut{
		Question: story.Question,
		Round:    story.Round,
		StoryAnswersInOut: StoryAnswersInOut{
			FibbingIt: answers,
		},
	}
	return newStory, nil
}

func (f FibbingIt) newAnswersOut(storyAnswers *FibbingItAnswers) FibbingItAnswersInOut {
	var answers FibbingItAnswersInOut
	for _, storyAnswer := range *storyAnswers {
		answer := FibbingItAnswerInOut(storyAnswer)
		answers = append(answers, answer)
	}

	return answers
}

func (f FibbingIt) NewStory(story StoryInOut) (Story, error) {
	answers, err := f.newAnswers(story.StoryAnswersInOut.FibbingIt)
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

func (f FibbingIt) newAnswers(storyAnswers FibbingItAnswersInOut) (FibbingItAnswers, error) {
	var answers FibbingItAnswers
	if len(storyAnswers) == 0 {
		return FibbingItAnswers{}, errors.BadRequestf("no answers in the story.")
	}

	for _, storyAnswer := range storyAnswers {
		answer := FibbingItAnswer(storyAnswer)
		answers = append(answers, answer)
	}

	return answers, nil
}
