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
	"banter-bus-server/tests/data"
)

func (s *Tests) SubTestAddGame(t *testing.T) {
	for _, tc := range data.AddGame {
		testName := fmt.Sprintf("Add New Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
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
	for _, tc := range data.GetGame {
		testName := fmt.Sprintf("Get Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
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
	for _, tc := range data.RemoveGame {
		testName := fmt.Sprintf("Remove Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
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
	for _, tc := range data.EnableGame {
		testName := fmt.Sprintf("Enable Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			enableOrDisableTest(t, s.router, "enable", tc.Name, tc.ExpectedGameNames, tc.ExpectedStatus)
		})
	}
}

func (s *Tests) SubTestDisableGame(t *testing.T) {
	for _, tc := range data.DisableGame {
		testName := fmt.Sprintf("Disable Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
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
