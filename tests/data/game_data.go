package data

import (
	"net/http"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// AddGame is the test data for adding new game types.
var AddGame = []struct {
	TestDescription string
	Payload         interface{}
	ExpectedStatus  int
	ExpectedGame    serverModels.Game
}{
	{
		"Add a new game",
		&serverModels.NewGame{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		},
		http.StatusCreated,
		serverModels.Game{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
			Enabled:  true,
		},
	},
	{
		"Try to add another game wrong Nam field",
		struct{ Nam, RulesURL string }{
			Nam:      "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		}, http.StatusBadRequest,
		serverModels.Game{},
	},
	{
		"Try to add another game wrong Rule field",
		struct{ Name, RulURL string }{
			Name:   "fibbing_it",
			RulURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		},
		http.StatusBadRequest,
		serverModels.Game{},
	},
	{
		"Try to add a game that already exists.",
		&serverModels.NewGame{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		}, http.StatusConflict,
		serverModels.Game{},
	},
}

// GetAllGames is the test data for getting existing game types.
var GetAllGames = []struct {
	TestDescription string
	Filter          string
	ExpectedNames   []string
}{
	{
		"Get games no filter",
		"",
		[]string{
			"quibly",
			"fibbing_it",
			"drawlosseum",
		},
	},
	{
		"Get games all filter",
		"all",
		[]string{
			"quibly",
			"fibbing_it",
			"drawlosseum",
		},
	},
	{
		"Get games enabled filter",
		"enabled",
		[]string{
			"quibly",
			"fibbing_it",
		},
	},
	{
		"Get games disabled filter",
		"disabled",
		[]string{
			"drawlosseum",
		},
	},
}

// GetGame is the test data for getting existing game types.
var GetGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    serverModels.Game
}{
	{
		"Get a game",
		"quibly",
		http.StatusOK,
		serverModels.Game{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
			Enabled:  true,
		},
	},
	{
		"Get another game",
		"fibbing_it",
		http.StatusOK,
		serverModels.Game{
			Name:     "fibbing_it",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbing_it",
			Enabled:  true,
		},
	},
	{
		"Try to get game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		serverModels.Game{},
	},
	{
		"Try to get another game that doesn't exist",
		"another_one",
		http.StatusNotFound,
		serverModels.Game{},
	},
}

// RemoveGame is the test data for removing existing game types.
var RemoveGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
}{
	{
		"Remove an existing game",
		"quibly",
		http.StatusOK,
	},
	{
		"Try to remove a game that's already been removed",
		"quibly",
		http.StatusNotFound,
	},
	{
		"Try to remove another game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
	},
}

// EnableGame is the test data for enabling game types.
var EnableGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    serverModels.Game
}{
	{
		"Enable a disabled game",
		"drawlosseum",
		http.StatusOK,
		serverModels.Game{
			Name:     "drawlosseum",
			RulesURL: "https://google.com/drawlosseum",
			Enabled:  true,
		},
	},
	{
		"Enable an already enabled game",
		"drawlosseum",
		http.StatusOK,
		serverModels.Game{
			Name:     "drawlosseum",
			RulesURL: "https://google.com/drawlosseum",
			Enabled:  true,
		},
	},
	{
		"Enable a game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		serverModels.Game{},
	},
}

// DisableGame is the test data for disabling game types.
var DisableGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    serverModels.Game
}{
	{
		"Disable an enabled game",
		"fibbing_it",
		http.StatusOK,
		serverModels.Game{
			Name:     "fibbing_it",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbing_it",
			Enabled:  false,
		},
	},
	{
		"Disable an already disabled game",
		"fibbing_it",
		http.StatusOK,
		serverModels.Game{
			Name:     "fibbing_it",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbing_it",
			Enabled:  false,
		},
	},
	{
		"Disable a game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		serverModels.Game{},
	},
}
