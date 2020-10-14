package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"banter-bus-server/src/core/config"
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server"
	"banter-bus-server/src/server/models"

	"github.com/wI2L/fizz"
	"gopkg.in/go-playground/assert.v1"
)

var router *fizz.Fizz

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	database.RemoveCollection("game")
	os.Exit(code)
}

func setup() {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "config.test.yml")
	config := config.GetConfig()
	fmt.Println(config.Database.Host)
	dbConfig := database.DatabaseConfig{
		Username:     config.Database.Username,
		Password:     config.Database.Password,
		DatabaseName: config.Database.DatabaseName,
		Host:         config.Database.Host,
		Port:         config.Database.Port,
	}
	database.InitialiseDatabase(dbConfig)
	router, _ = server.NewRouter()
}

func TestCreateGame(t *testing.T) {
	cases := []struct {
		Payload  interface{}
		Expected int
	}{
		{&models.NewGame{
			Name:     "quibly",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
		}, http.StatusOK},
		{&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusOK},
		{struct{ Nam, RulesURL string }{
			Nam:      "quiblyv3",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest},
		{struct{ Name, RuleURL string }{
			Name:    "quiblyv3",
			RuleURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest},
		{struct{ Nam, RuleURL string }{
			Nam:     "quiblyv3",
			RuleURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusBadRequest},
		{&models.NewGame{
			Name:     "quiblyv2",
			RulesURL: "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
		}, http.StatusConflict},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Add New Game"), func(t *testing.T) {
			data, _ := json.Marshal(tc.Payload)
			encodedData := bytes.NewBuffer([]byte(data))
			req, _ := http.NewRequest("POST", "/game", encodedData)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.Expected, w.Code)
		})
	}
}

func TestGetAllGames(t *testing.T) {
	req, _ := http.NewRequest("GET", "/game", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var response []string
	json.Unmarshal([]byte(w.Body.String()), &response)

	var expectedResult = []string{
		"quibly",
		"quiblyv2",
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResult, response)
}
