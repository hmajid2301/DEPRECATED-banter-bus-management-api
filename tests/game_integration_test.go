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

func (s *Tests) SubTestAddGame(t *testing.T) {
	cases := []struct {
		Payload  interface{}
		Expected int
	}{
		{&models.NewGame{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		}, http.StatusOK},
		{&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusOK},
		{struct{ Nam, RulesURL string }{
			Nam:      "quiblyv3",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest},
		{struct{ Name, RuleURL string }{
			Name:    "quiblyv3",
			RuleURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest},
		{struct{ Nam, RuleURL string }{
			Nam:     "quiblyv3",
			RuleURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest},
		{&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusConflict},
	}

	for _, tc := range cases {
		t.Run("Add New Game", func(t *testing.T) {
			data, _ := json.Marshal(tc.Payload)
			encodedData := bytes.NewBuffer(data)
			req, _ := http.NewRequest("POST", "/game", encodedData)
			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			assert.Equal(t, tc.Expected, w.Code)
		})
	}
}

func (s *Tests) SubTestGetAllGames(t *testing.T) {
	req, _ := http.NewRequest("GET", "/game", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	var response []string
	err := json.Unmarshal(w.Body.Bytes(), &response)

	if err != nil {
		fmt.Println(err)
	}

	var expectedResult = []string{
		"a_game",
		"fibbly",
		"draw_me",
		"new_totally_original_game",
		"new_totally_original_game_2",
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResult, response)
}

func (s *Tests) SubTestGetGame(t *testing.T) {
	cases := []struct {
		Name           string
		ExpectedStatus int
		ExpectedGame   models.Game
	}{
		{
			"a_game",
			http.StatusOK,
			models.Game{
				Name:      "a_game",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/a_game",
				Questions: &models.Question{},
				Enabled:   true,
			},
		},
		{
			"fibbly",
			http.StatusOK,
			models.Game{
				Name:      "fibbly",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbly",
				Questions: &models.Question{},
				Enabled:   true,
			},
		},
		{
			"quiblyv3",
			http.StatusNotFound,
			models.Game{},
		},
		{
			"another_one",
			http.StatusNotFound,
			models.Game{},
		},
	}

	for _, tc := range cases {
		t.Run("Get Game", func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/game/%s", tc.Name), nil)
			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			assert.Equal(t, tc.ExpectedStatus, w.Code)

			if w.Code == http.StatusOK {
				var response *models.Game

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					fmt.Println(err)
				}

				assert.Equal(t, tc.ExpectedGame, response)
			}
		})
	}
}

func (s *Tests) SubTestRemoveGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames []string
	}{
		{
			"a_game",
			http.StatusOK,
			[]string{"fibbly", "draw_me", "new_totally_original_game", "new_totally_original_game_2"},
		},
		{
			"a_game",
			http.StatusNotFound,
			[]string{},
		},
		{
			"quiblyv3",
			http.StatusNotFound,
			[]string{},
		},
		{
			"another_one",
			http.StatusNotFound,
			[]string{},
		},
	}

	for _, tc := range cases {
		t.Run("Remove Game", func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/game/%s", tc.Name), nil)
			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			assert.Equal(t, tc.ExpectedStatus, w.Code)

			if w.Code == http.StatusOK {
				req, _ := http.NewRequest("GET", "/game", nil)
				w := httptest.NewRecorder()
				s.router.ServeHTTP(w, req)
				var response []string

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					fmt.Println(err)
				}

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, tc.ExpectedGameNames, response)
			}
		})
	}
}

func (s *Tests) SubTestEnableGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames models.Game
	}{
		{
			"draw_me",
			http.StatusOK,
			models.Game{
				Name:      "draw_me",
				RulesURL:  "https://google.com/draw_me",
				Questions: &models.Question{},
				Enabled:   true,
			},
		},
		{
			"draw_me",
			http.StatusConflict,
			models.Game{},
		},
		{
			"quiblyv2",
			http.StatusNotFound,
			models.Game{},
		},
	}

	for _, tc := range cases {
		t.Run("Enable A Game", func(t *testing.T) {
			enableOrDisableTest(t, s.router, "enable", tc.Name, tc.ExpectedGameNames, tc.ExpectedStatus)
		})
	}
}

func (s *Tests) SubTestDisableGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames models.Game
	}{
		{
			"a_game",
			http.StatusOK,
			models.Game{
				Name:      "a_game",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/a_game",
				Questions: &models.Question{},
				Enabled:   false,
			},
		},
		{
			"a_game",
			http.StatusConflict,
			models.Game{},
		},
		{
			"quiblyv2",
			http.StatusNotFound,
			models.Game{},
		},
	}

	for _, tc := range cases {
		t.Run("Disable A Game", func(t *testing.T) {
			enableOrDisableTest(t, s.router, "disable", tc.Name, tc.ExpectedGameNames, tc.ExpectedStatus)
		})
	}
}

func enableOrDisableTest(
	t *testing.T,
	router http.Handler,
	enable string,
	gameName string,
	expectedGameNames models.Game,
	expectedStatus int) {
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/game/%s/%s", gameName, enable), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, expectedStatus, w.Code)

	if w.Code == http.StatusOK {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/game/%s", gameName), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var response *models.Game
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedGameNames, response)
	}
}
