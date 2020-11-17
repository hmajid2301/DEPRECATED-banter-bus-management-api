package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"banter-bus-server/tests/data"

	"gopkg.in/go-playground/assert.v1"
)

func (s *Tests) SubTestHealthcheck(t *testing.T) {
	for _, tc := range data.Healthcheck {
		testName := fmt.Sprintf("Healthcheck: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/healthcheck", nil)
			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			assert.Equal(t, tc.Expected, w.Code)
		})
	}
}
