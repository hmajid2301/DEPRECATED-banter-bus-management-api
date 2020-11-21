package data

import (
	"net/http"

	"banter-bus-server/src/server/models"
)

// AddGame is the test data for adding new game types.
var AddGame = []struct {
	TestDescription string
	Payload         interface{}
	ExpectedStatus  int
	ExpectedGame    models.Game
}{
	{
		"Add a new game",
		&models.NewGame{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		},
		http.StatusCreated,
		models.Game{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
			Enabled:  true,
		},
	},
	{
		"Add another new game",
		&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		},
		http.StatusCreated,
		models.Game{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
			Enabled:  true,
		},
	},
	{
		"Try to add another game wrong Nam field",
		struct{ Nam, RulesURL string }{
			Nam:      "quiblyv3",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest,
		models.Game{},
	},
	{
		"Try to add another game wrong Rule field",
		struct{ Name, RuleURL string }{
			Name:    "quiblyv3",
			RuleURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		},
		http.StatusBadRequest,
		models.Game{},
	},
	{
		"Try to add a game that already exists.",
		&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusConflict,
		models.Game{},
	},
}

// GetAllGame is the test data for getting existing game types.
var GetAllGame = []struct {
	TestDescription string
	Filter          string
	ExpectedNames   []string
}{
	{
		"Get games no filter",
		"",
		[]string{
			"a_game",
			"fibbly",
			"draw_me",
			"new_totally_original_game",
			"new_totally_original_game_2",
		},
	},
	{
		"Get games all filter",
		"all",
		[]string{
			"a_game",
			"fibbly",
			"draw_me",
			"new_totally_original_game",
			"new_totally_original_game_2",
		},
	},
	{
		"Get games enabled filter",
		"enabled",
		[]string{
			"a_game",
			"fibbly",
			"new_totally_original_game",
		},
	},
	{
		"Get games disabled filter",
		"disabled",
		[]string{
			"draw_me",
			"new_totally_original_game_2",
		},
	},
}

// GetGame is the test data for getting existing game types.
var GetGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    models.Game
}{
	{
		"Get a game",
		"a_game",
		http.StatusOK,
		models.Game{
			Name:     "a_game",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/a_game",
			Enabled:  true,
		},
	},
	{
		"Get another game",
		"fibbly",
		http.StatusOK,
		models.Game{
			Name:     "fibbly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbly",
			Enabled:  true,
		},
	},
	{
		"Try to get game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		models.Game{},
	},
	{
		"Try to get another game that doesn't exist",
		"another_one",
		http.StatusNotFound,
		models.Game{},
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
		"a_game",
		http.StatusOK,
	},
	{
		"Try to remove a game that's already been removed",
		"a_game",
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
	ExpectedGame    models.Game
}{
	{
		"Enabled a disabled game",
		"draw_me",
		http.StatusOK,
		models.Game{
			Name:     "draw_me",
			RulesURL: "https://google.com/draw_me",
			Enabled:  true,
		},
	},
	{
		"Enabled an already enabled game",
		"draw_me",
		http.StatusConflict,
		models.Game{},
	},
	{
		"Enabled a game that doesn't exist",
		"quiblyv2",
		http.StatusNotFound,
		models.Game{},
	},
}

// DisableGame is the test data for disabling game types.
var DisableGame = []struct {
	TestDescription string
	Name            string
	ExpectedStatus  int
	ExpectedGame    models.Game
}{
	{
		"Disable an enabled game",
		"a_game",
		http.StatusOK,
		models.Game{
			Name:     "a_game",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/a_game",
			Enabled:  false,
		},
	},
	{
		"Disable an already disabled game",
		"a_game",
		http.StatusConflict,
		models.Game{},
	},
	{
		"Disable a game that doesn't exist",
		"quiblyv2",
		http.StatusNotFound,
		models.Game{},
	},
}
