package controllers_test

import (
	"fmt"
	"testing"

	"banter-bus-server/tests/data"
)

func (s *Tests) SubTestAddQuestionToGame(t *testing.T) {
	for _, tc := range data.AddQuestion {
		testName := fmt.Sprintf("Add New Question: %s", tc.TestDescription)
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
	// Add some questions
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
