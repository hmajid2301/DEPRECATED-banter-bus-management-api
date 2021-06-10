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
			endpoint := fmt.Sprintf("/game/%s/question", tc.Game)
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
			endpoint := fmt.Sprintf("/game/%s/question", tc.Game)
			s.httpExpect.DELETE(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
			// TODO: test question cannot be found after
		})
	}
}

func (s *Tests) SubTestAddTranslationQuestion(t *testing.T) {
	for _, tc := range data.AddTranslationQuestion {
		testName := fmt.Sprintf("Add Question Translation: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/%s", tc.Game, tc.LanguageCode)
			s.httpExpect.POST(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestRemoveTranslation(t *testing.T) {
	for _, tc := range data.RemoveTranslationQuestion {
		testName := fmt.Sprintf("Remove Question Translation: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/%s", tc.Game, tc.LanguageCode)
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
			endpoint := fmt.Sprintf("/game/%s/question/enable", tc.Game)
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
			endpoint := fmt.Sprintf("/game/%s/question/disable", tc.Game)
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
			endpoint := fmt.Sprintf("/game/%s/question/group", tc.Payload.Name)
			response := s.httpExpect.GET(endpoint).WithJSON(tc.Payload).WithQueryObject(tc.Payload)
			retval := response.Expect().Status(tc.ExpectedCode)

			if tc.ExpectedCode == http.StatusOK {
				retval.JSON().Array().Equal(tc.ExpectedGroups)
			}
		})
	}
}
