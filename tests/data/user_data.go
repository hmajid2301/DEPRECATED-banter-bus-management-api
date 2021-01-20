package data

import (
	"net/http"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

var admin = true
var notAdmin = !admin

//AddUser is the test data for adding new users.
var AddUser = []struct {
	TestDescription string
	Payload         interface{}
	ExpectedStatus  int
}{
	{
		"Add a new user",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "paid",
			Admin:      &admin,
		},
		http.StatusCreated,
	},

	{
		"Add a user that has already been added",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "paid",
			Admin:      &admin,
		},
		http.StatusConflict,
	},

	{
		"Try to add a user with a taken username but different admin privileges",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "paid",
			Admin:      &notAdmin,
		},
		http.StatusConflict,
	},

	{
		"Try to add a user with a taken username but different membership",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "free",
			Admin:      &admin,
		},
		http.StatusConflict,
	},

	{
		"Try to add a user with an invalid membership",
		&serverModels.NewUser{
			Username:   "seeb123",
			Membership: "overlord",
			Admin:      &notAdmin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with no username",
		&serverModels.NewUser{
			Username:   "",
			Membership: "paid",
			Admin:      &admin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with invalid fields 1",
		struct {
			Usernam, Membership string
			Admin               *bool
		}{
			Usernam:    "seeb123",
			Membership: "free",
			Admin:      &notAdmin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with invalid fields 2",
		struct {
			Username, Membershap string
			Admin                *bool
		}{
			Username:   "seeb123",
			Membershap: "free",
			Admin:      &notAdmin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with invalid fields 3",
		struct {
			Username, Membership string
			Admine               *bool
		}{
			Username:   "seeb123",
			Membership: "free",
			Admine:     &notAdmin,
		},
		http.StatusBadRequest,
	},
}

//GetAllUsers is the test data for geting all users
var GetAllUsers = []struct {
	TestDescription   string
	Filter            *serverModels.ListUserParams
	ExpectedUsernames []string
}{
	{
		"Get all users no filter",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "all",
			Membership:  "all",
		},
		[]string{
			"virat_kohli",
			"roh1t_sharma",
			"dhawanShikhar",
		},
	},
	{
		"Get all admins",
		&serverModels.ListUserParams{
			AdminStatus: "admin",
			Privacy:     "all",
			Membership:  "all",
		},
		[]string{
			"roh1t_sharma",
		},
	},
	{
		"Get all non-admins",
		&serverModels.ListUserParams{
			AdminStatus: "non_admin",
			Privacy:     "all",
			Membership:  "all",
		},
		[]string{
			"virat_kohli",
			"dhawanShikhar",
		},
	},
	{
		"Get all members with public profiles",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "public",
			Membership:  "all",
		},
		[]string{
			"virat_kohli",
		},
	},
	{
		"Get all members with private profiles",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "private",
			Membership:  "all",
		},
		[]string{
			"roh1t_sharma",
			"dhawanShikhar",
		},
	},
	{
		"Get all members with free membership",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "all",
			Membership:  "free",
		},
		[]string{
			"roh1t_sharma",
		},
	},
	{
		"Get all members with paid membership",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "all",
			Membership:  "paid",
		},
		[]string{
			"virat_kohli",
			"dhawanShikhar",
		},
	},
	{
		"Get members who are admin, have a private account and free membership",
		&serverModels.ListUserParams{
			AdminStatus: "admin",
			Privacy:     "private",
			Membership:  "free",
		},
		[]string{
			"roh1t_sharma",
		},
	},
}

//GetUser is the test data for getting users.
var GetUser = []struct {
	TestDescription string
	Username        string
	ExpectedStatus  int
	ExpectedUser    serverModels.User
}{
	{
		"Get a user",
		"virat_kohli",
		http.StatusOK,
		serverModels.User{
			Username:   "virat_kohli",
			Admin:      &notAdmin,
			Privacy:    "public",
			Membership: "paid",
			Preferences: &serverModels.UserPreferences{
				LanguageCode: "pa",
			},
			Friends: []serverModels.Friend{
				{
					Username: "roh1t_sharma",
				},
				{
					Username: "dhawanShikhar",
				},
			},
		},
	},
	{
		"Get another user",
		"roh1t_sharma",
		http.StatusOK,
		serverModels.User{
			Username:   "roh1t_sharma",
			Admin:      &admin,
			Privacy:    "private",
			Membership: "free",
			Preferences: &serverModels.UserPreferences{
				LanguageCode: "mr",
			},
			Friends: []serverModels.Friend{
				{
					Username: "virat_kohli",
				},
			},
		},
	},
	{
		"Get a user that doesn't exists",
		"azharAli",
		http.StatusNotFound,
		serverModels.User{},
	},
}

// GetUserPools is the test data for getting users.
var GetUserPools = []struct {
	TestDescription string
	Username        string
	ExpectedStatus  int
	ExpectedResult  []serverModels.QuestionPool
}{
	{
		"Get user pool for a user",
		"virat_kohli",
		http.StatusOK,
		[]serverModels.QuestionPool{
			{
				PoolName:     "my_pool",
				GameName:     "fibbing_it",
				LanguageCode: "en",
				Privacy:      "public",
				Questions: serverModels.QuestionPoolQuestions{
					FibbingIt: serverModels.FibbingItQuestionsPool{
						Likely: []string{
							"to eat ice-cream from the tub",
							"to get arrested",
						},
						FreeForm: map[string][]string{
							"bike_group": {
								"Favourite bike colour?",
								"A funny question?",
							},
						},
						Opinion: map[string]map[string][]string{
							"horse_group": {
								"questions": []string{
									"What do you think about horses?",
									"What do you think about camels?",
								},
								"answers": []string{"lame", "tasty"},
							},
						},
					},
				},
			},
			{
				PoolName:     "my_pool2",
				GameName:     "quibly",
				LanguageCode: "fr",
				Privacy:      "private",
				Questions: serverModels.QuestionPoolQuestions{
					Quibly: serverModels.QuiblyQuestionsPool{
						Pair: []string{
							"What do you think about horses?",
							"What do you think about camels?",
						},
						Answers: []string{
							"Favourite bike colour?",
							"A funny question?",
						},
					},
				},
			},
		},
	},
	{
		"Get another user pool",
		"dhawanShikhar",
		http.StatusOK,
		[]serverModels.QuestionPool{
			{
				PoolName:     "draw_me",
				GameName:     "drawlosseum",
				LanguageCode: "en",
				Privacy:      "friends",
				Questions: serverModels.QuestionPoolQuestions{
					Drawlosseum: serverModels.DrawlosseumQuestionsPool{
						Drawings: []string{
							"horses",
							"camels",
						},
					},
				},
			},
			{
				PoolName:     "my_unique_pool2",
				GameName:     "quibly",
				LanguageCode: "en",
				Privacy:      "public",
				Questions: serverModels.QuestionPoolQuestions{
					Quibly: serverModels.QuiblyQuestionsPool{
						Group: []string{
							"What do you think about horses?",
							"What do you think about camels?",
						},
					},
				},
			},
		},
	},
	{
		"Get a user pool for a user that doesn't exists",
		"azharAli",
		http.StatusNotFound,
		[]serverModels.QuestionPool{},
	},
}

// GetSingleUserPool is the data for testing a single pool returned from a user.
var GetSingleUserPool = []struct {
	TestDescription string
	Username        string
	PoolName        string
	ExpectedStatus  int
	ExpectedResult  serverModels.QuestionPool
}{
	{
		"Get user pool for a user",
		"virat_kohli",
		"my_pool",
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			LanguageCode: "en",
			Privacy:      "public",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
					},
					FreeForm: map[string][]string{
						"bike_group": {
							"Favourite bike colour?",
							"A funny question?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about horses?",
								"What do you think about camels?",
							},
							"answers": []string{"lame", "tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Get another user pool",
		"dhawanShikhar",
		"draw_me",
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "draw_me",
			GameName:     "drawlosseum",
			LanguageCode: "en",
			Privacy:      "friends",
			Questions: serverModels.QuestionPoolQuestions{
				Drawlosseum: serverModels.DrawlosseumQuestionsPool{
					Drawings: []string{
						"horses",
						"camels",
					},
				},
			},
		},
	},
	{
		"Get a user pool that doesn't exist",
		"dhawanShikhar",
		"draw_me1",
		http.StatusNotFound,
		serverModels.QuestionPool{},
	},
	{
		"Get a user pool for a user that doesn't exists",
		"azharAli",
		"a_pool",
		http.StatusNotFound,
		serverModels.QuestionPool{},
	},
}

// RemoveUser is the data for testing removing users
var RemoveUser = []struct {
	TestDescription string
	Username        string
	ExpectedStatus  int
}{
	{
		"Remove an existing user",
		"virat_kohli",
		http.StatusOK,
	},
	{
		"Remove a user that's already been removed",
		"virat_kohli",
		http.StatusNotFound,
	},
	{
		"Remove another existing user",
		"roh1t_sharma",
		http.StatusOK,
	},
	{
		"Try to remove a non-existent user",
		"NaseemShah",
		http.StatusNotFound,
	},
}

// GetUserStories is the test data for getting user's stories.
var GetUserStories = []struct {
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

// AddNewPool is the data for testing adding new pools to users.
var AddNewPool = []struct {
	TestDescription string
	Username        string
	NewPool         interface{}
	ExpectedStatus  int
	ExpectedResult  serverModels.QuestionPool
}{
	{
		"Add new user pool",
		"virat_kohli",
		serverModels.NewQuestionPool{
			PoolName: "another_pool",
			GameName: "fibbing_it",
			Privacy:  "public",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "another_pool",
			GameName:     "fibbing_it",
			Privacy:      "public",
			LanguageCode: "en",
		},
	},
	{
		"Add new user pool to another user",
		"roh1t_sharma",
		serverModels.NewQuestionPool{
			PoolName:     "first_pool",
			GameName:     "fibbing_it",
			LanguageCode: "fr",
			Privacy:      "private",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "first_pool",
			GameName:     "fibbing_it",
			LanguageCode: "fr",
			Privacy:      "private",
		},
	},
	{
		"Add new user pool, wrong privacy setting",
		"virat_kohli",
		serverModels.NewQuestionPool{
			PoolName: "another_pool2",
			GameName: "fibbing_it",
			Privacy:  "wrong",
		},
		http.StatusBadRequest,
		serverModels.QuestionPool{},
	},
	{
		"Add new user pool, wrong language code",
		"virat_kohli",
		serverModels.NewQuestionPool{
			PoolName:     "another_pool2",
			GameName:     "fibbing_it",
			LanguageCode: "papa",
			Privacy:      "private",
		},
		http.StatusBadRequest,
		serverModels.QuestionPool{},
	},
	{
		"Add new user pool, wrong privacy setting",
		"virat_kohli",
		serverModels.NewQuestionPool{
			PoolName: "another_pool2",
			GameName: "fibbing_it",
			Privacy:  "wrong",
		},
		http.StatusBadRequest,
		serverModels.QuestionPool{},
	},
	{
		"Add new user pool, game doesn't exist",
		"virat_kohli",
		serverModels.NewQuestionPool{
			PoolName: "another_pool2",
			GameName: "fibbing2_it",
			Privacy:  "public",
		},
		http.StatusBadRequest,
		serverModels.QuestionPool{},
	},
	{
		"Add new user pool incorrect field PoolName",
		"virat_kohli",
		struct{ PoolNam, GameName, Privacy string }{
			PoolNam:  "another_pool2",
			GameName: "fibbing_it",
			Privacy:  "public",
		},
		http.StatusBadRequest,
		serverModels.QuestionPool{},
	},
	{
		"Add new user pool incorrect field PoolName",
		"virat_kohli",
		struct{ PoolName, GameNam, Privacy string }{
			PoolName: "another_pool4",
			GameNam:  "fibbing_it",
			Privacy:  "public",
		},
		http.StatusBadRequest,
		serverModels.QuestionPool{},
	},
	{
		"Add new user pool, user doesn't exist",
		"virat",
		serverModels.NewQuestionPool{
			PoolName: "another_pool4",
			GameName: "fibbing_it",
			Privacy:  "public",
		},
		http.StatusNotFound,
		serverModels.QuestionPool{},
	},
	{
		"Add existing pool",
		"virat_kohli",
		serverModels.NewQuestionPool{
			PoolName: "another_pool",
			GameName: "fibbing_it",
			Privacy:  "public",
		},
		http.StatusConflict,
		serverModels.QuestionPool{},
	},
}

// RemovePool is the data for testing removing pools from users.
var RemovePool = []struct {
	TestDescription string
	Username        string
	PoolName        interface{}
	ExpectedStatus  int
}{
	{
		"Remove user pool",
		"virat_kohli",
		"my_pool",
		http.StatusOK,
	},
	{
		"Remove pool from another user",
		"dhawanShikhar",
		"draw_me",
		http.StatusOK,
	},
	{
		"Remove user pool, user doesn't exist",
		"virat",
		"another_pool4",
		http.StatusNotFound,
	},
	{
		"Remove an already deleted pool",
		"virat_kohli",
		"another_pool",
		http.StatusNotFound,
	},
}

// UpdatePool is the data for testing adding new pools to users.
var UpdatePool = []struct {
	TestDescription string
	Username        string
	PoolName        string
	UpdatePool      interface{}
	ExpectedStatus  int
	ExpectedResult  serverModels.QuestionPool
}{
	{
		"Add question to user pool",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "do you love bikes?",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "horse_group",
					Type: "questions",
				},
			},
			Operation: "add",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			LanguageCode: "en",
			GameName:     "fibbing_it",
			Privacy:      "public",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
					},
					FreeForm: map[string][]string{
						"bike_group": {
							"Favourite bike colour?",
							"A funny question?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about horses?",
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Add another question to user pool",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "super tasty",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "horse_group",
					Type: "answers",
				},
			},
			Operation: "add",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			LanguageCode: "en",
			Privacy:      "public",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
					},
					FreeForm: map[string][]string{
						"bike_group": {
							"Favourite bike colour?",
							"A funny question?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about horses?",
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty", "super tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Add question to user pool different round",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "what is the tastiest horse?",
				Round:   "free_form",
				Group: &serverModels.Group{
					Name: "bike_group",
				},
			},
			Operation: "add",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			LanguageCode: "en",
			Privacy:      "public",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
					},
					FreeForm: map[string][]string{
						"bike_group": {
							"Favourite bike colour?",
							"A funny question?",
							"what is the tastiest horse?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about horses?",
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty", "super tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Add question to user pool another different round",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "to eat a fruit",
				Round:   "likely",
				Group:   nil,
			},
			Operation: "add",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			LanguageCode: "en",
			Privacy:      "public",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
						"to eat a fruit",
					},
					FreeForm: map[string][]string{
						"bike_group": {
							"Favourite bike colour?",
							"A funny question?",
							"what is the tastiest horse?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about horses?",
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty", "super tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Add question to another user pool",
		"virat_kohli",
		"my_pool2",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "what do you think about donkeys?",
				Round:   "pair",
				Group:   nil,
			},
			Operation: "add",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool2",
			GameName:     "quibly",
			LanguageCode: "fr",
			Privacy:      "private",
			Questions: serverModels.QuestionPoolQuestions{
				Quibly: serverModels.QuiblyQuestionsPool{
					Pair: []string{
						"What do you think about horses?",
						"What do you think about camels?",
						"what do you think about donkeys?",
					},
					Answers: []string{
						"Favourite bike colour?",
						"A funny question?",
					},
					Group: []string{},
				},
			},
		},
	},
	{
		"Add another question to another user pool",
		"virat_kohli",
		"my_pool2",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "a question?",
				Round:   "group",
				Group:   nil,
			},
			Operation: "add",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool2",
			GameName:     "quibly",
			LanguageCode: "fr",
			Privacy:      "private",
			Questions: serverModels.QuestionPoolQuestions{
				Quibly: serverModels.QuiblyQuestionsPool{
					Pair: []string{
						"What do you think about horses?",
						"What do you think about camels?",
						"what do you think about donkeys?",
					},
					Answers: []string{
						"Favourite bike colour?",
						"A funny question?",
					},
					Group: []string{
						"a question?",
					},
				},
			},
		},
	},
	{
		"Add another question to another user theit pool",
		"dhawanShikhar",
		"draw_me",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "chicken",
			},
			Operation: "add",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "draw_me",
			GameName:     "drawlosseum",
			LanguageCode: "en",
			Privacy:      "friends",
			Questions: serverModels.QuestionPoolQuestions{
				Drawlosseum: serverModels.DrawlosseumQuestionsPool{
					Drawings: []string{
						"horses",
						"camels",
						"chicken",
					},
				},
			},
		},
	},
	{
		"Remove question from user pool",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "What do you think about horses?",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "horse_group",
					Type: "questions",
				},
			},
			Operation: "remove",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			LanguageCode: "en",
			Privacy:      "public",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
						"to eat a fruit",
					},
					FreeForm: map[string][]string{
						"bike_group": []string{
							"Favourite bike colour?",
							"A funny question?",
							"what is the tastiest horse?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty", "super tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Remove another question from user pool",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "super tasty",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "horse_group",
					Type: "answers",
				},
			},
			Operation: "remove",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			Privacy:      "public",
			LanguageCode: "en",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
						"to eat a fruit",
					},
					FreeForm: map[string][]string{
						"bike_group": {
							"Favourite bike colour?",
							"A funny question?",
							"what is the tastiest horse?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Remove question from user pool different round",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "what is the tastiest horse?",
				Round:   "free_form",
				Group: &serverModels.Group{
					Name: "bike_group",
				},
			},
			Operation: "remove",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			Privacy:      "public",
			LanguageCode: "en",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
						"to eat a fruit",
					},
					FreeForm: map[string][]string{
						"bike_group": {
							"Favourite bike colour?",
							"A funny question?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Remove question from user pool another different round",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "to eat a fruit",
				Round:   "likely",
				Group:   nil,
			},
			Operation: "remove",
		},
		http.StatusOK,
		serverModels.QuestionPool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			Privacy:      "public",
			LanguageCode: "en",
			Questions: serverModels.QuestionPoolQuestions{
				FibbingIt: serverModels.FibbingItQuestionsPool{
					Likely: []string{
						"to eat ice-cream from the tub",
						"to get arrested",
					},
					FreeForm: map[string][]string{
						"bike_group": []string{
							"Favourite bike colour?",
							"A funny question?",
						},
					},
					Opinion: map[string]map[string][]string{
						"horse_group": {
							"questions": []string{
								"What do you think about camels?",
								"do you love bikes?",
							},
							"answers": []string{"lame", "tasty"},
						},
					},
				},
			},
		},
	},
	{
		"Remove question from user pool that was already removed",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "to eat a fruit",
				Round:   "likely",
				Group:   nil,
			},
			Operation: "remove",
		},
		http.StatusNotFound,
		serverModels.QuestionPool{},
	},
	{
		"Remove another question from user pool that was already removed",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "what is the tastiest horse?",
				Round:   "free_form",
				Group: &serverModels.Group{
					Name: "bike_group",
				},
			},
			Operation: "remove",
		},
		http.StatusNotFound,
		serverModels.QuestionPool{},
	},
	{
		"Add question to user pool that already exists",
		"virat_kohli",
		"my_pool",
		serverModels.UpdateQuestionPool{
			NewQuestion: serverModels.NewQuestion{
				Content: "to eat ice-cream from the tub",
				Round:   "likely",
				Group:   nil,
			},
			Operation: "add",
		},
		http.StatusConflict,
		serverModels.QuestionPool{},
	},
}
