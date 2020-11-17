package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"banter-bus-server/tests/data"
)

func (s *Tests) SubTestAddGame(t *testing.T) {
	for _, tc := range data.AddGame {
		testName := fmt.Sprintf("Add New Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			s.httpExpect.POST("/game").
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestGetAllGames(t *testing.T) {
	var expectedResult = []string{
		"a_game",
		"fibbly",
		"draw_me",
		"new_totally_original_game",
		"new_totally_original_game_2",
	}
	s.httpExpect.GET("/game").
		Expect().
		Status(http.StatusOK).JSON().Array().Equal(expectedResult)
}

func (s *Tests) SubTestGetGame(t *testing.T) {
	for _, tc := range data.GetGame {
		testName := fmt.Sprintf("Get Game: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s", tc.Name)
			response := s.httpExpect.GET(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				response.JSON().Object().Equal(tc.ExpectedGame)
			}
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
				s.httpExpect.GET("/game").
					Expect().
					Status(http.StatusOK).JSON().Array().Equal(tc.ExpectedGameNames)
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
		})
	}
}
