package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/tests/data"

	"github.com/gavv/httpexpect"
)

func (s *Tests) SubTestAddUser(t *testing.T) {
	for _, tc := range data.AddUser {
		testName := fmt.Sprintf("Add New User: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			s.httpExpect.POST("/user").
				WithJSON(tc.Payload).
				Expect().
				Status(tc.ExpectedStatus)
		})
	}
}

func (s *Tests) SubTestGetUser(t *testing.T) {
	for _, tc := range data.GetUser {
		testName := fmt.Sprintf("Get User: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			getUser(tc.Username, tc.ExpectedStatus, tc.ExpectedUser, s.httpExpect)
		})
	}
}

func (s *Tests) SubTestGetAllUsers(t *testing.T) {
	for _, tc := range data.GetAllUsers {
		testName := fmt.Sprintf("Get All Users: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			response := s.httpExpect.GET("/user").WithQueryObject(tc.Filter)
			response.
				Expect().
				Status(http.StatusOK).JSON().Array().Equal(tc.ExpectedUsernames)
		})
	}
}

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

func (s *Tests) SubTestRemoveUser(t *testing.T) {
	for _, tc := range data.RemoveUser {
		testName := fmt.Sprintf("Remove User: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s", tc.Username)
			s.httpExpect.DELETE(endpoint).Expect().Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				getUser(tc.Username, http.StatusNotFound, serverModels.User{}, s.httpExpect)
			}
		})
	}
}

func (s *Tests) SubTestGetUserStory(t *testing.T) {
	for _, tc := range data.GetUserStories {
		testName := fmt.Sprintf("Get User Stories: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s/story", tc.Username)
			response := s.httpExpect.GET(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

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
				newPool, ok := tc.NewPool.(serverModels.NewQuestionPool)
				if !ok {
					t.Errorf("failed to convert to NewQuestionPool")
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

			if tc.ExpectedStatus == http.StatusOK {
				endpoint := fmt.Sprintf("/user/%s/pool/%s", tc.Username, tc.PoolName)
				s.httpExpect.GET(endpoint).
					Expect().
					Status(http.StatusNotFound)
			}
		})
	}
}

func (s *Tests) SubTestUpdateUserPool(t *testing.T) {
	for _, tc := range data.UpdatePool {
		testName := fmt.Sprintf("Update Pool: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/user/%s/pool/%s", tc.Username, tc.PoolName)
			s.httpExpect.PATCH(endpoint).
				WithJSON(tc.UpdatePool).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				response := getUserPool(tc.Username, tc.PoolName, s.httpExpect, http.StatusOK)
				response.JSON().Object().Equal(tc.ExpectedResult)
			}
		})
	}
}

func getUser(user string, expectedStatus int, expectedResult serverModels.User, httpExpect *httpexpect.Expect) {
	endpoint := fmt.Sprintf("/user/%s", user)
	response := httpExpect.GET(endpoint).
		Expect().
		Status(expectedStatus)

	if expectedStatus == http.StatusOK {
		response.JSON().Object().Equal(expectedResult)
	}
}

func getAllUserPools(username string, httpExpect *httpexpect.Expect, expectedStatus int) *httpexpect.Response {
	endpoint := fmt.Sprintf("/user/%s/pool", username)
	response := httpExpect.GET(endpoint).
		Expect().
		Status(expectedStatus)

	return response
}

func getUserPool(
	username string,
	poolName string,
	httpExpect *httpexpect.Expect,
	expectedStatus int,
) *httpexpect.Response {
	endpoint := fmt.Sprintf("/user/%s/pool/%s", username, poolName)
	response := httpExpect.GET(endpoint).
		Expect().
		Status(expectedStatus)

	return response
}
