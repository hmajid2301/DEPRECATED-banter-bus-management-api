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
