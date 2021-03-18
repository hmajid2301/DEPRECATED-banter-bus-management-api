package data

import (
	"net/http"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// GetUserPools is the test data for getting users.
var GetUserPools = []struct {
	TestDescription string
	Username        string
	ExpectedStatus  int
	ExpectedResult  []serverModels.Pool
}{
	{
		"Get user pool for a user",
		"virat_kohli",
		http.StatusOK,
		[]serverModels.Pool{
			{
				PoolName:     "my_pool",
				GameName:     "fibbing_it",
				LanguageCode: "en",
				Privacy:      "public",
			},
			{
				PoolName:     "my_pool2",
				GameName:     "quibly",
				LanguageCode: "fr",
				Privacy:      "private",
			},
		},
	},
	{
		"Get another user pool",
		"dhawanShikhar",
		http.StatusOK,
		[]serverModels.Pool{
			{
				PoolName:     "draw_me",
				GameName:     "drawlosseum",
				LanguageCode: "en",
				Privacy:      "friends",
			},
		},
	},
	{
		"Get a user pool for a user that doesn't exists",
		"azharAli",
		http.StatusNotFound,
		[]serverModels.Pool{},
	},
}

// GetSingleUserPool is the data for testing a single pool returned from a user.
var GetSingleUserPool = []struct {
	TestDescription string
	Username        string
	PoolName        string
	ExpectedStatus  int
	ExpectedResult  serverModels.Pool
}{
	{
		"Get user pool for a user",
		"virat_kohli",
		"my_pool",
		http.StatusOK,
		serverModels.Pool{
			PoolName:     "my_pool",
			GameName:     "fibbing_it",
			LanguageCode: "en",
			Privacy:      "public",
		},
	},
	{
		"Get another user pool",
		"dhawanShikhar",
		"draw_me",
		http.StatusOK,
		serverModels.Pool{
			PoolName:     "draw_me",
			GameName:     "drawlosseum",
			LanguageCode: "en",
			Privacy:      "friends",
		},
	},
	{
		"Get a user pool that doesn't exist",
		"dhawanShikhar",
		"draw_me1",
		http.StatusNotFound,
		serverModels.Pool{},
	},
	{
		"Get a user pool for a user that doesn't exists",
		"azharAli",
		"a_pool",
		http.StatusNotFound,
		serverModels.Pool{},
	},
}

// AddNewPool is the data for testing adding new pools to users.
var AddNewPool = []struct {
	TestDescription string
	Username        string
	NewPool         interface{}
	ExpectedStatus  int
	ExpectedResult  serverModels.Pool
}{
	{
		"Add new user pool",
		"virat_kohli",
		serverModels.Pool{
			PoolName: "another_pool",
			GameName: "fibbing_it",
			Privacy:  "public",
		},
		http.StatusOK,
		serverModels.Pool{
			PoolName:     "another_pool",
			GameName:     "fibbing_it",
			Privacy:      "public",
			LanguageCode: "en",
		},
	},
	{
		"Add new user pool to another user",
		"roh1t_sharma",
		serverModels.Pool{
			PoolName:     "first_pool",
			GameName:     "fibbing_it",
			LanguageCode: "fr",
			Privacy:      "private",
		},
		http.StatusOK,
		serverModels.Pool{
			PoolName:     "first_pool",
			GameName:     "fibbing_it",
			LanguageCode: "fr",
			Privacy:      "private",
		},
	},
	{
		"Add new user pool, wrong privacy setting",
		"virat_kohli",
		serverModels.Pool{
			PoolName: "another_pool2",
			GameName: "fibbing_it",
			Privacy:  "wrong",
		},
		http.StatusBadRequest,
		serverModels.Pool{},
	},
	{
		"Add new user pool, wrong language code",
		"virat_kohli",
		serverModels.Pool{
			PoolName:     "another_pool2",
			GameName:     "fibbing_it",
			LanguageCode: "papa",
			Privacy:      "private",
		},
		http.StatusBadRequest,
		serverModels.Pool{},
	},
	{
		"Add new user pool, wrong privacy setting",
		"virat_kohli",
		serverModels.Pool{
			PoolName: "another_pool2",
			GameName: "fibbing_it",
			Privacy:  "wrong",
		},
		http.StatusBadRequest,
		serverModels.Pool{},
	},
	{
		"Add new user pool, game doesn't exist",
		"virat_kohli",
		serverModels.Pool{
			PoolName: "another_pool2",
			GameName: "fibbing2_it",
			Privacy:  "public",
		},
		http.StatusBadRequest,
		serverModels.Pool{},
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
		serverModels.Pool{},
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
		serverModels.Pool{},
	},
	{
		"Add new user pool, user doesn't exist",
		"virat",
		serverModels.Pool{
			PoolName: "another_pool4",
			GameName: "fibbing_it",
			Privacy:  "public",
		},
		http.StatusNotFound,
		serverModels.Pool{},
	},
	{
		"Add existing pool",
		"virat_kohli",
		serverModels.Pool{
			PoolName: "another_pool",
			GameName: "fibbing_it",
			Privacy:  "public",
		},
		http.StatusConflict,
		serverModels.Pool{},
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
		"",
		http.StatusNotFound,
	},
	{
		"Remove user pool, pool doesn't exist (wildcard)",
		"virat",
		"*",
		http.StatusNotFound,
	},
	{
		"Remove user pool, invalid pool",
		"virat_kohli",
		"",
		http.StatusNotFound,
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

// AddQuestionPool is the data for testing adding new pools to users.
var AddQuestionPool = []struct {
	TestDescription string
	Username        string
	PoolName        string
	UpdatePool      interface{}
	ExpectedStatus  int
	ExpectedResult  serverModels.Pool
}{
	{
		"Add question to user pool",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "do you love bikes?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "question",
			},
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Add another question to user pool",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "super tasty",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "answer",
			},
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Add question to user pool different round",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "what is the tastiest horse?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Add question to user pool another different round",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "to eat a fruit",
			Round:   "likely",
			Group:   nil,
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Add question to another user pool",
		"virat_kohli",
		"my_pool2",
		serverModels.NewQuestion{
			Content: "what do you think about donkeys?",
			Round:   "pair",
			Group:   nil,
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Add another question to another user pool",
		"virat_kohli",
		"my_pool2",
		serverModels.NewQuestion{
			Content: "a question?",
			Round:   "group",
			Group:   nil,
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Add another question to another user theit pool",
		"dhawanShikhar",
		"draw_me",
		serverModels.NewQuestion{
			Content: "chicken",
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Add question to user that doesn't exist",
		"azhar",
		"my_pool",
		serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
			Group:   nil,
		},
		http.StatusNotFound,
		serverModels.Pool{},
	},
	{
		"Add question to pool that doesn't exist",
		"virat_kohli",
		"my",
		serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
			Group:   nil,
		},
		http.StatusNotFound,
		serverModels.Pool{},
	},
	{
		"Add question to user pool that already exists",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
			Group:   nil,
		},
		http.StatusConflict,
		serverModels.Pool{},
	},
}

// RemoveQuestionPool is the data for testing removing question from pools.
var RemoveQuestionPool = []struct {
	TestDescription string
	Username        string
	PoolName        string
	UpdatePool      interface{}
	ExpectedStatus  int
	ExpectedResult  serverModels.Pool
}{
	{
		"Remove question from user pool",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "What do you think about horses?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "question",
			},
		},
		http.StatusOK,
		serverModels.Pool{},
	},
	{
		"Remove question from user pool that was already removed",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "to eat a fruit",
			Round:   "likely",
			Group:   nil,
		},
		http.StatusNotFound,
		serverModels.Pool{},
	},
	{
		"Remove another question from user pool that was already removed",
		"virat_kohli",
		"my_pool",
		serverModels.NewQuestion{
			Content: "what is the tastiest horse?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		},
		http.StatusNotFound,
		serverModels.Pool{},
	},
	{
		"Remove question to user that doesn't exist",
		"azhar",
		"my_pool",
		serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
			Group:   nil,
		},
		http.StatusNotFound,
		serverModels.Pool{},
	},
	{
		"Remove question to pool that doesn't exist",
		"virat_kohli",
		"my",
		serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
			Group:   nil,
		},
		http.StatusNotFound,
		serverModels.Pool{},
	},
}
