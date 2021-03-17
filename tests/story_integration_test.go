package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/tests/data"
)

func (s *Tests) SubTestGetStories(t *testing.T) {
	for _, tc := range data.GetStories {
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
