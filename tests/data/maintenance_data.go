package data

import (
	"net/http"
)

// Healthcheck is the test data for add questions to a game types.
var Healthcheck = []struct {
	TestDescription string
	Expected        int
}{
	{
		"Healthcheck",
		http.StatusOK,
	},
}
