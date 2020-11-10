package data

import (
	"banter-bus-server/src/server/models"
	"net/http"
)

// AddQuestion is the test data for add questions to a game types.
var AddQuestion = []struct {
	TestDescription string
	GameType        string
	Payload         interface{}
	Expected        int
}{
	{
		"Add a question to round one",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "what is the funniest thing ever told?",
			Round:   "one",
		}, http.StatusOK,
	},
	{
		"Add another question to round one",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "This is another question?",
			Round:   "one",
		}, http.StatusOK,
	},
	{
		"Add question to round two",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "what is the funniest thing ever told?",
			Round:   "two",
		}, http.StatusOK,
	},
	{
		"Add question to round three",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "what is the funniest thing ever told?",
			Round:   "three",
		}, http.StatusOK,
	},
	{
		"Try add question to wrong round (four)",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "what is the funniest thing ever told?",
			Round:   "four",
		}, http.StatusBadRequest,
	},
	{
		"Try add question wrong Question field",
		"new_totally_original_game",
		struct{ Que, Round string }{
			Que:   "quibly",
			Round: "one",
		}, http.StatusBadRequest,
	},
	{
		"Try add question wrong Round field",
		"new_totally_original_game",
		struct{ Question, Rod string }{
			Question: "quibly",
			Rod:      "one",
		}, http.StatusBadRequest,
	},
	{
		"Try to add question to not existent game",
		"does_not_exist",
		&models.NewQuestion{
			Content: "What is a question?",
			Round:   "one",
		}, http.StatusNotFound,
	},
	{
		"Try to add question that already exists",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "what is the funniest thing ever told?",
			Round:   "three",
		}, http.StatusConflict,
	},
	{
		"Try to add question that already exists round one",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "this is a question?",
			Round:   "one",
		}, http.StatusConflict,
	},
	{
		"Try to add another question that already exists round three",
		"new_totally_original_game",
		&models.NewQuestion{
			Content: "this is a another question?",
			Round:   "three",
		}, http.StatusConflict,
	},
	{
		"Try to add questions that already exists round one to another game type",
		"new_totally_original_game_2",
		&models.NewQuestion{
			Content: "what is the funniest thing ever told?",
			Round:   "one",
		}, http.StatusConflict,
	},
}
