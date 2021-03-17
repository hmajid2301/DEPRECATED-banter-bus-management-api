package data

import (
	"net/http"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// GetStories is the test data for getting user's stories.
var GetStories = []struct {
	TestDescription string
	Username        string
	ExpectedStatus  int
	ExpectedResult  []serverModels.Story
}{
	{
		"Get user's stories for a user",
		"virat_kohli",
		http.StatusOK,
		[]serverModels.Story{
			{
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
			{
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
			{
				Question: "What do you think about horses?",
				Round:    "opinion",
				StoryAnswers: serverModels.StoryAnswers{
					FibbingIt: []serverModels.StoryFibbingIt{
						{
							Nickname: "!sus",
							Answer:   "tasty",
						},
						{
							Nickname: "normal_guy1",
							Answer:   "lame",
						},
						{
							Nickname: "normal_girl1",
							Answer:   "lame",
						},
						{
							Nickname: "normal_person1",
							Answer:   "lame",
						},
					},
				},
			},
		},
	},
	{
		"Get another user's stories",
		"roh1t_sharma",
		http.StatusOK,
		[]serverModels.Story{
			{
				Question: "what do you think about horses?",
				Round:    "free_form",
				StoryAnswers: serverModels.StoryAnswers{
					FibbingIt: []serverModels.StoryFibbingIt{
						{
							Nickname: "!sus",
							Answer:   "tasty",
						},
						{
							Nickname: "normal_guy1",
							Answer:   "hello",
						},
						{
							Nickname: "normal_girl1",
							Answer:   "what is a horse?",
						},
						{
							Nickname: "normal_person1",
							Answer:   "is this a real game?",
						},
					},
				},
			},
			{
				Question: "most likely to get arrested?",
				Round:    "likely",
				StoryAnswers: serverModels.StoryAnswers{
					FibbingIt: []serverModels.StoryFibbingIt{
						{Answer: "normal_guy1", Nickname: "!sus"},
						{Answer: "normal_girl1", Nickname: "normal_guy1"},
						{Answer: "!sus", Nickname: "normal_girl1"},
						{Answer: "normal_girl1", Nickname: "normal_person1"},
					},
				},
			},
		},
	},
	{
		"Get another user's stories, empty",
		"dhawanShikhar",
		http.StatusOK,
		[]serverModels.Story{},
	},
	{
		"Get a user's stories for a user that doesn't exists",
		"azharAli",
		http.StatusNotFound,
		[]serverModels.Story{},
	},
}
