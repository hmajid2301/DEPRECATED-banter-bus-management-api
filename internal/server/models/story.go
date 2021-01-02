package serverModels

// Story struct to contain information about a user story
type Story struct {
	Question string `json:"question"`
	Round    string `json:"round,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	StoryAnswers
}

// StoryAnswers contains all the different stories answers that are supported, for different game types.
type StoryAnswers struct {
	Drawlosseum []StoryDrawlosseum `json:"drawlosseum,omitempty"`
	Quibly      []StoryQuibly      `json:"quibly,omitempty"`
	FibbingIt   []StoryFibbingIt   `json:"fibbing_it,omitempty"`
}

// StoryFibbingIt contains information about the Fibbing It answers for user stories
type StoryFibbingIt struct {
	Nickname string `json:"nickname,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

// StoryQuibly contains information about the Quibly answers for user stories
type StoryQuibly struct {
	Nickname string `json:"nickname,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Votes    int    `json:"votes,omitempty"`
}

// StoryDrawlosseum contains information about the Drawlosseum answers for user stories
type StoryDrawlosseum struct {
	Start DrawlosseumDrawingPoint `json:"start,omitempty"`
	End   DrawlosseumDrawingPoint `json:"end,omitempty"`
	Color string                  `json:"color,omitempty"`
}

// DrawlosseumDrawingPoint contains information about a point in a Drawlosseum drawing
type DrawlosseumDrawingPoint struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}
