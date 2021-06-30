package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/tests/data"
)

func (s *Tests) SubTestAddStories(t *testing.T) {
	for _, tc := range data.AddStories {
		testName := fmt.Sprintf("Add Story: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/story/%s", tc.GameName)
			s.httpExpect.POST(endpoint).
				WithJSON(tc.Payload).
				Expect().
				Status(tc.ExpectedStatus)
		})
	}
}

func (s *Tests) SubTestGetStories(t *testing.T) {
	for _, tc := range data.GetStories {
		testName := fmt.Sprintf("Get Story: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/story/%s", tc.StoryID)
			response := s.httpExpect.GET(endpoint).
				Expect().
				Status(tc.ExpectedStatus)

			if tc.ExpectedStatus == http.StatusOK {
				response.JSON().Equal(tc.ExpectedResult)
			}
		})
	}
}

func (s *Tests) SubTestDeleteStories(t *testing.T) {
	for _, tc := range data.GetStories {
		testName := fmt.Sprintf("Delete Story: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			endpoint := fmt.Sprintf("/story/%s", tc.StoryID)
			s.httpExpect.DELETE(endpoint).
				Expect().
				Status(tc.ExpectedStatus)
		})
	}
}
