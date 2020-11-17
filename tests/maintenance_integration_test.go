package controllers_test

import (
	"fmt"
	"testing"

	"banter-bus-server/tests/data"
)

func (s *Tests) SubTestHealthcheck(t *testing.T) {
	for _, tc := range data.Healthcheck {
		testName := fmt.Sprintf("Healthcheck: %s", tc.TestDescription)
		t.Run(testName, func(t *testing.T) {
			s.httpExpect.GET("/healthcheck").
				Expect().
				Status(tc.Expected)
		})
	}
}
