package data

import (
	"net/http"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/games"
)

var AddGame = []struct {
	TestDescription string
	Payload         interface{}
	ExpectedStatus  int
	ExpectedGame    games.GameOut
}{
	{
		"Add a new game",
		&games.GameIn{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		},
		http.StatusCreated,
		games.GameOut{
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
		games.GameOut{},
	},
	{
		"Try to add another game wrong Rule field",
		struct{ Name, RulURL string }{
			Name:   "fibbing_it",
			RulURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		},
		http.StatusBadRequest,
		games.GameOut{},
	},
	{
		"Try to add a game that already exists.",
		&games.GameIn{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		}, http.StatusConflict,
		games.GameOut{},
	},
}

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

var GetGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    games.GameOut
}{
	{
		"Get a game",
		"quibly",
		http.StatusOK,
		games.GameOut{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
			Enabled:  true,
		},
	},
	{
		"Get another game",
		"fibbing_it",
		http.StatusOK,
		games.GameOut{
			Name:     "fibbing_it",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbing_it",
			Enabled:  true,
		},
	},
	{
		"Try to get game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		games.GameOut{},
	},
	{
		"Try to get another game that doesn't exist",
		"another_one",
		http.StatusNotFound,
		games.GameOut{},
	},
}

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

var EnableGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    games.GameOut
}{
	{
		"Enable a disabled game",
		"drawlosseum",
		http.StatusOK,
		games.GameOut{
			Name:     "drawlosseum",
			RulesURL: "https://google.com/drawlosseum",
			Enabled:  true,
		},
	},
	{
		"Enable an already enabled game",
		"drawlosseum",
		http.StatusOK,
		games.GameOut{
			Name:     "drawlosseum",
			RulesURL: "https://google.com/drawlosseum",
			Enabled:  true,
		},
	},
	{
		"Enable a game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		games.GameOut{},
	},
}

var DisableGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    games.GameOut
}{
	{
		"Disable an enabled game",
		"fibbing_it",
		http.StatusOK,
		games.GameOut{
			Name:     "fibbing_it",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbing_it",
			Enabled:  false,
		},
	},
	{
		"Disable an already disabled game",
		"fibbing_it",
		http.StatusOK,
		games.GameOut{
			Name:     "fibbing_it",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbing_it",
			Enabled:  false,
		},
	},
	{
		"Disable a game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		games.GameOut{},
	},
}
