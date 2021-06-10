package data

import (
	"net/http"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// GetStories is the test data for getting user's stories.
var GetStories = []struct {
	TestDescription string
	StoryID         string
	ExpectedStatus  int
	ExpectedResult  serverModels.Story
}{
	{
		"Get a story",
		"1def4233-f674-4a3f-863d-6e850bfbfdb4",
		http.StatusOK,
		serverModels.Story{
			Question: "how many fish are there?",
			Round:    "pair",
			StoryAnswers: serverModels.StoryAnswers{
				Quibly: []serverModels.StoryQuibly{
					{
						Nickname: "funnyMan420",
						Answer:   "one",
						Votes:    12341,
					},
					{
						Nickname: "123456",
						Answer:   "many",
						Votes:    0,
					},
				},
			},
		},
	},
	{
		"Get another story",
		"a4ffd1c8-93c5-4f4c-8ace-71996edcbcb7",
		http.StatusOK,
		serverModels.Story{
			Question: "fish",
			Nickname: "i_cannotDraw",
			StoryAnswers: serverModels.StoryAnswers{
				Drawlosseum: []serverModels.StoryDrawlosseum{
					{
						Start: serverModels.DrawlosseumDrawingPoint{
							X: 100,
							Y: -100,
						},
						End: serverModels.DrawlosseumDrawingPoint{
							X: 90,
							Y: -100,
						},
						Color: "#000",
					},
				},
			},
		},
	},
	{
		"Story does not exist",
		"50-011c-45d8-98f7-819520c253b6",
		http.StatusNotFound,
		serverModels.Story{},
	},
}
