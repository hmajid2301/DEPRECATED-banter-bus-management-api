package controllers_test

import (
	"fmt"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/tests/data"
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
