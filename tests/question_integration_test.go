package controllers_test

import (
	"banter-bus-server/tests/data"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func (s *Tests) SubTestAddQuestionToGame(t *testing.T) {
	for _, tc := range data.AddQuestion {
		testName := fmt.Sprintf("Add New Question: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			data, _ := json.Marshal(tc.Payload)
			encodedData := bytes.NewBuffer(data)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/game/%s/question", tc.GameType), encodedData)
			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			assert.Equal(t, tc.Expected, w.Code)
		})
	}
}
