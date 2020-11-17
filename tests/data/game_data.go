package data

import (
	"net/http"

	"banter-bus-server/src/server/models"
)

// AddGame is the test data for adding new game types.
var AddGame = []struct {
	TestDescription string
	Payload         interface{}
	Expected        int
}{
	{
		"Add a new game",
		&models.NewGame{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		},
		http.StatusCreated,
	},
	{
		"Add another new game",
		&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		},
		http.StatusCreated,
	},
	{
		"Try to add another game wrong Nam field",
		struct{ Nam, RulesURL string }{
			Nam:      "quiblyv3",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest,
	},
	{
		"Try to add another game wrong Rule field",
		struct{ Name, RuleURL string }{
			Name:    "quiblyv3",
			RuleURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		},
		http.StatusBadRequest,
	},
	{
		"Try to add a game that already exists.",
		&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusConflict,
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
	TestDescription   string
	Name              string
	ExpectedStatus    int
	ExpectedGameNames []string
}{
	{
		"Remove an existing game",
		"a_game",
		http.StatusOK,
		[]string{"fibbly", "draw_me", "new_totally_original_game", "new_totally_original_game_2"},
	},
	{
		"Try to remove a game thats already been removed",
		"a_game",
		http.StatusNotFound,
		[]string{},
	},
	{
		"Try to remove another game that doesn't exist",
		"quiblyv3",
		http.StatusNotFound,
		[]string{},
	},
}

// EnableGame is the test data for enabling game types.
var EnableGame = []struct {
	TestDescription   string
	Name              string
	ExpectedStatus    int
	ExpectedGameNames models.Game
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
	TestDescription   string
	Name              string
	ExpectedStatus    int
	ExpectedGameNames models.Game
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
