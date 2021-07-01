package story

import (
	"github.com/juju/errors"
)

type Drawlosseum struct{}

func (d Drawlosseum) NewStoryOut(story Story) (StoryInOut, error) {
	storyAnswers, ok := story.Answers.(*DrawlosseumAnswers)
	if !ok {
		return StoryInOut{}, errors.Errorf("invalid answer for Drawlosseum")
	}
	answers := d.newAnswersOut(storyAnswers)
	newStory := StoryInOut{
		Question: story.Question,
		Nickname: story.Nickname,
		StoryAnswersInOut: StoryAnswersInOut{
			Drawlosseum: answers,
		},
	}
	return newStory, nil
}

func (d Drawlosseum) newAnswersOut(storyAnswers *DrawlosseumAnswers) DrawlosseumAnswersInOut {
	var answers DrawlosseumAnswersInOut
	for _, storyAnswer := range *storyAnswers {
		answers = append(answers, storyAnswer)
	}

	return answers
}

func (d Drawlosseum) NewStory(story StoryInOut) (Story, error) {
	answers, err := d.newAnswers(story.StoryAnswersInOut.Drawlosseum)
	if err != nil {
		return Story{}, err
	}

	newStory := Story{
		Question: story.Question,
		Nickname: story.Nickname,
		Answers:  answers,
	}
	return newStory, nil
}

func (d Drawlosseum) newAnswers(storyAnswers DrawlosseumAnswersInOut) (DrawlosseumAnswers, error) {
	var answers DrawlosseumAnswers
	if len(storyAnswers) == 0 {
		return DrawlosseumAnswers{}, errors.BadRequestf("no answers in the story.")
	}

	for _, storyAnswer := range storyAnswers {
		answers = append(answers, storyAnswer)
	}

	return answers, nil
}
