package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/tests/data"
)

func (s *Tests) SubTestGetUserPools(t *testing.T) {
	for _, tc := range data.GetUserPools {
		testName := fmt.Sprintf("Get All User Pools: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			response := getAllUserPools(tc.Username, s.httpExpect, tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				response.JSON().Array().Equal(tc.ExpectedResult)
			}
		})
	}
}

func (s *Tests) SubTestAddUserPool(t *testing.T) {
	for _, tc := range data.AddNewPool {
		testName := fmt.Sprintf("Add New Pool: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s/pool", tc.Username)
			s.httpExpect.POST(endpoint).
				WithJSON(tc.NewPool).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				newPool, ok := tc.NewPool.(serverModels.Pool)
				if !ok {
					t.Errorf("failed to convert to QuestionPool")
				}

				endpoint := fmt.Sprintf("/user/%s/pool/%s", tc.Username, newPool.PoolName)
				s.httpExpect.GET(endpoint).
					Expect().
					Status(tc.ExpectedStatus).JSON().Equal(tc.ExpectedResult)
			}
		})
	}
}

func (s *Tests) SubTestGetUserPool(t *testing.T) {
	for _, tc := range data.GetSingleUserPool {
		testName := fmt.Sprintf("Get Single User Pool: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s/pool/%s", tc.Username, tc.PoolName)
			response := s.httpExpect.GET(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				response.JSON().Equal(tc.ExpectedResult)
			}
		})
	}
}

func (s *Tests) SubTestRemoveUserPool(t *testing.T) {
	for _, tc := range data.RemovePool {
		testName := fmt.Sprintf("Remove Pool: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s/pool/%s", tc.Username, tc.PoolName)
			s.httpExpect.DELETE(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

			// TODO: test 0 questions in question collection related to this pool. After #29 is merged in.
			if tc.ExpectedStatus == http.StatusOK {
				endpoint := fmt.Sprintf("/user/%s/pool/%s", tc.Username, tc.PoolName)
				s.httpExpect.GET(endpoint).
					Expect().
					Status(http.StatusNotFound)
			}
		})
	}
}

func (s *Tests) SubTestAddQuestionToPool(t *testing.T) {
	for _, tc := range data.AddQuestionPool {
		testName := fmt.Sprintf("Add Question to Pool: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s/pool/%s/question", tc.Username, tc.PoolName)
			s.httpExpect.POST(endpoint).
				WithJSON(tc.UpdatePool).
				Expect().
				Status(tc.ExpectedStatus)

			// TODO: test questions is in pool after #29
		})
	}
}

func (s *Tests) SubTestRemoveQuestionFromPool(t *testing.T) {
	for _, tc := range data.RemoveQuestionPool {
		testName := fmt.Sprintf("Remove Question from Pool: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s/pool/%s/question", tc.Username, tc.PoolName)
			s.httpExpect.DELETE(endpoint).
				WithJSON(tc.UpdatePool).
				Expect().
				Status(tc.ExpectedStatus)

			// TODO: test questions is in pool after #29
		})
	}
}

func getAllUserPools(username string, httpExpect *httpexpect.Expect, expectedStatus int) *httpexpect.Response {
	endpoint := fmt.Sprintf("/user/%s/pool", username)
	response := httpExpect.GET(endpoint).
		Expect().
		Status(expectedStatus)

	return response
}
