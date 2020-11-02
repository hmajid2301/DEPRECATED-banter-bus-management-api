package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"

	"banter-bus-server/src/server/models"
)

func (s *Tests) SubTestAddQuestionToGame(t *testing.T) {
	cases := []struct {
		GameType string
		Payload  interface{}
		Expected int
	}{
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "what is the funniest thing ever told?",
				Round:   "one",
			}, http.StatusOK,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "This is another question?",
				Round:   "one",
			}, http.StatusOK,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "what is the funniest thing ever told?",
				Round:   "two",
			}, http.StatusOK,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "what is the funniest thing ever told?",
				Round:   "three",
			}, http.StatusOK,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "this is a question?",
				Round:   "two",
			}, http.StatusOK,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "what is the funniest thing ever told?",
				Round:   "four",
			}, http.StatusBadRequest,
		},
		{
			"new_totally_original_game",
			struct{ Que, Round string }{
				Que:   "quibly",
				Round: "one",
			}, http.StatusBadRequest,
		},
		{
			"new_totally_original_game",
			struct{ Question, Rod string }{
				Question: "quibly",
				Rod:      "one",
			}, http.StatusBadRequest,
		},
		{
			"does_not_exist",
			&models.NewQuestion{
				Content: "What is a question?",
				Round:   "one",
			}, http.StatusNotFound,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "what is the funniest thing ever told?",
				Round:   "three",
			}, http.StatusConflict,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "this is a question?",
				Round:   "one",
			}, http.StatusConflict,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "this is a another question?",
				Round:   "three",
			}, http.StatusConflict,
		},
		{
			"new_totally_original_game",
			&models.NewQuestion{
				Content: "what is the funniest thing ever told?",
				Round:   "one",
			}, http.StatusConflict,
		},
		{
			"new_totally_original_game_2",
			&models.NewQuestion{
				Content: "what is the funniest thing ever told?",
				Round:   "one",
			}, http.StatusConflict,
		},
	}

	for _, tc := range cases {
		t.Run("Add New Question", func(t *testing.T) {
			data, _ := json.Marshal(tc.Payload)
			encodedData := bytes.NewBuffer(data)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/game/%s/question", tc.GameType), encodedData)
			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			assert.Equal(t, tc.Expected, w.Code)
		})
	}
}
