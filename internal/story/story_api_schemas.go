package story

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal"
)

type StoryInOut struct {
	Question string `json:"question"`
	Round    string `json:"round,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	StoryAnswersInOut
}

type StoryAnswersInOut struct {
	Drawlosseum DrawlosseumAnswersInOut `json:"drawlosseum,omitempty"`
	Quibly      QuiblyAnswersInOut      `json:"quibly,omitempty"`
	FibbingIt   FibbingItAnswersInOut   `json:"fibbing_it,omitempty"`
}

func (s StoryAnswersInOut) Answer() {}

type FibbingItAnswerInOut struct {
	Nickname string `json:"nickname"`
	Answer   string `json:"answer"`
}

type FibbingItAnswersInOut []FibbingItAnswerInOut

type QuiblyAnswerInOut struct {
	Nickname string `json:"nickname"`
	Answer   string `json:"answer"`
	Votes    int    `json:"votes"`
}

type QuiblyAnswersInOut []QuiblyAnswerInOut

type DrawlosseumAnswersInOut []CaertsianCoordinateColor

type StoryIDParams struct {
	StoryID string `description:"The id for the story." example:"2b45f6c6-d8be-4d13-9fc6-2f821c925774" path:"story_id"`
}

type CurrentStoryInput struct {
	internal.GameParams
	StoryIDParams
}

type NewStoryInput struct {
	internal.GameParams
	StoryInOut
}
