package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

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

func (s *Tests) SubTestGetQuestions(t *testing.T) {
	for _, tc := range data.GetQuestions {
		testName := fmt.Sprintf("Get Questions: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question", tc.Game)
			response := s.httpExpect.GET(endpoint).
				WithQuery("round", tc.Round).WithQuery("language", tc.Language).WithQuery("limit", tc.Limit).
				WithQuery("group_name", tc.GroupName).WithQuery("enabled", tc.Enabled).WithQuery("random", tc.Random).
				Expect().Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				if tc.Random == true {
					response.JSON().Array().Length().Equal(tc.Limit)
				} else {
					response.JSON().Array().Equal(tc.ExpectedQuestions)
				}
			}
		})
	}
}

func (s *Tests) SubTestGetQuestionsIds(t *testing.T) {
	for _, tc := range data.GetAllQuestionsIds {
		testName := fmt.Sprintf("Get All Question IDs: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			response := s.getQuestionsByID(tc.Game, tc.Limit, tc.Cursor, tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				response.JSON().Equal(tc.ExpectedPayload)
			}
		})
	}
}

func (s *Tests) SubTestRemoveQuestionFromGame(t *testing.T) {
	for _, tc := range data.RemoveQuestion {
		testName := fmt.Sprintf("Remove Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/%s", tc.Game, tc.ID)
			s.httpExpect.DELETE(endpoint).
				Expect().
				Status(tc.Expected)

			if tc.Expected == http.StatusOK {
				getQuestionEndpoint := fmt.Sprintf("%s/en", endpoint)
				s.httpExpect.GET(getQuestionEndpoint).
					Expect().
					Status(http.StatusNotFound)
			}
		})
	}
}

func (s *Tests) SubTestAddTranslationQuestion(t *testing.T) {
	for _, tc := range data.AddTranslationQuestion {
		testName := fmt.Sprintf("Add Question Translation: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/%s/%s", tc.Game, tc.ID, tc.LanguageCode)
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
			endpoint := fmt.Sprintf("/game/%s/question/%s/%s", tc.Game, tc.ID, tc.LanguageCode)
			s.httpExpect.DELETE(endpoint).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestGetQuestionByID(t *testing.T) {
	for _, tc := range data.GetQuestionById {
		testName := fmt.Sprintf("Get Question By ID: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/%s/%s", tc.Game, tc.ID, tc.LanguageCode)
			response := s.httpExpect.GET(endpoint).
				Expect().
				Status(tc.Expected)

			if tc.Expected == http.StatusOK {
				response.JSON().Equal(tc.ExpectedPayload)
			}
		})
	}
}

func (s *Tests) SubTestEnableQuestion(t *testing.T) {
	for _, tc := range data.EnableQuestion {
		testName := fmt.Sprintf("Enable Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/%s/enable", tc.Game, tc.ID)
			s.httpExpect.PUT(endpoint).
				Expect().
				Status(tc.Expected)
		})
	}
}

func (s *Tests) SubTestDisableQuestion(t *testing.T) {
	for _, tc := range data.DisableQuestion {
		testName := fmt.Sprintf("Disable Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/game/%s/question/%s/disable", tc.Game, tc.ID)
			s.httpExpect.PUT(endpoint).
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
			response := s.httpExpect.GET(endpoint).WithQuery("round", tc.Payload.Round)
			retval := response.Expect().Status(tc.ExpectedCode)

			if tc.ExpectedCode == http.StatusOK {
				retval.JSON().Array().Equal(tc.ExpectedGroups)
			}
		})
	}
}

func (s *Tests) getQuestionsByID(game string, limit int64, cursor string, expectedStatus int) *httpexpect.Response {
	endpoint := fmt.Sprintf("/game/%s/question/id", game)
	response := s.httpExpect.GET(endpoint).WithQuery("limit", limit).WithQuery("cursor", cursor).
		Expect().
		Status(expectedStatus)

	return response
}
