package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/questions"
	"gitlab.com/banter-bus/banter-bus-management-api/tests/data"

	"github.com/gavv/httpexpect"
)

func (s *Tests) SubTestAddGame(t *testing.T) {
	for _, tc := range data.AddGame {
		testName := fmt.Sprintf("Add New Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			gameData, ok := tc.Payload.(*games.GameIn)

			if ok && tc.ExpectedStatus == http.StatusCreated {
				endpoint := fmt.Sprintf("/game/%s", gameData.Name)
				s.httpExpect.DELETE(endpoint).
					Expect().
					Status(http.StatusOK)
			}

			s.httpExpect.POST("/game").
				WithJSON(tc.Payload).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusCreated {
				s.getGame(tc.ExpectedGame.Name, http.StatusOK, tc.ExpectedGame)
			}
		})
	}
}

func (s *Tests) SubTestGetAllGames(t *testing.T) {
	for _, tc := range data.GetAllGames {
		testName := fmt.Sprintf("Get All Games: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			var response *httpexpect.Request

			if tc.Filter == "" {
				response = s.httpExpect.GET("/game")
			} else {
				response = s.httpExpect.GET("/game").WithQuery("games", tc.Filter)
			}

			response.
				Expect().
				Status(http.StatusOK).JSON().Array().Equal(tc.ExpectedNames)
		})
	}
}

func (s *Tests) SubTestGetGame(t *testing.T) {
	for _, tc := range data.GetGame {
		testName := fmt.Sprintf("Get Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			s.getGame(tc.Name, tc.ExpectedStatus, tc.ExpectedGame)
		})
	}
}

func (s *Tests) SubTestRemoveGame(t *testing.T) {
	for _, tc := range data.RemoveGame {
		testName := fmt.Sprintf("Remove Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s", tc.Name)
			s.httpExpect.DELETE(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				s.getGame(tc.Name, http.StatusNotFound, games.GameOut{})
				response := s.getQuestionsByID(tc.Name, 5, "", http.StatusOK)
				response.JSON().Object().Equal(questions.AllQuestionOut{
					IDs:    []string{},
					Cursor: "",
				})
			}
		})
	}
}

func (s *Tests) SubTestEnableGame(t *testing.T) {
	for _, tc := range data.EnableGame {
		testName := fmt.Sprintf("Enable Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/enable", tc.Name)
			s.httpExpect.PUT(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				s.getGame(tc.Name, http.StatusOK, tc.ExpectedGame)
			}
		})
	}
}

func (s *Tests) SubTestDisableGame(t *testing.T) {
	for _, tc := range data.DisableGame {
		testName := fmt.Sprintf("Disable Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/disable", tc.Name)
			s.httpExpect.PUT(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				s.getGame(tc.Name, http.StatusOK, tc.ExpectedGame)
			}
		})
	}
}

func (s *Tests) getGame(game string, expectedStatus int, expectedResult games.GameOut) {
	endpoint := fmt.Sprintf("/game/%s", game)
	response := s.httpExpect.GET(endpoint).
		Expect().
		Status(expectedStatus)

	if expectedStatus == http.StatusOK || expectedStatus == http.StatusCreated {
		response.JSON().Object().Equal(expectedResult)
	}
}
