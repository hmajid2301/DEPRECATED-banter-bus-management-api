package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/tests/data"
)

func (s *Tests) SubTestAddQuestionToGame(t *testing.T) {
	for _, tc := range data.AddQuestion {
		testName := fmt.Sprintf("Add Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question", tc.GameType)
			s.httpExpect.POST(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestRemoveQuestionFromGame(t *testing.T) {
	for _, tc := range data.RemoveQuestion {
		testName := fmt.Sprintf("Remove Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question", tc.GameType)
			s.httpExpect.DELETE(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestEnableQuestion(t *testing.T) {
	for _, tc := range data.EnableQuestion {
		testName := fmt.Sprintf("Enable Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/enable", tc.GameType)
			s.httpExpect.PUT(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestDisableQuestion(t *testing.T) {
	for _, tc := range data.DisableQuestion {
		testName := fmt.Sprintf("Disable Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/disable", tc.GameType)
			s.httpExpect.PUT(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestUpdateQuestion(t *testing.T) {
	for _, tc := range data.UpdateQuestion {
		testName := fmt.Sprintf("Update Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question", tc.GameType)
			s.httpExpect.PUT(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestGetAllGroups(t *testing.T) {
	for _, tc := range data.GetAllGroups {
		testName := fmt.Sprintf("Get All Groups: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/group", tc.Payload.GameName)
			response := s.httpExpect.GET(endpoint).WithJSON(tc.Payload).WithQueryObject(tc.Payload)
			retval := response.Expect().Status(tc.ExpectedCode)

			if tc.ExpectedCode == http.StatusOK {
				retval.JSON().Array().Equal(tc.ExpectedGroups)
			}
		})
	}
}
