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

func TestGetGame(t *testing.T) {
	cases := []struct {
		Name           string
		ExpectedStatus int
		ExpectedGame   models.Game
	}{
		{
			"quibly",
			http.StatusOK,
			models.Game{
				Name:      "quibly",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
				Questions: []models.Question{},
				Enabled:   true,
			},
		},
		{
			"quiblyv2",
			http.StatusOK,
			models.Game{
				Name:      "quiblyv2",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quiblyv2",
				Questions: []models.Question{},
				Enabled:   true,
			},
		},
		{
			"quiblyv3",
			http.StatusNotFound,
			models.Game{},
		},
		{
			"another_one",
			http.StatusNotFound,
			models.Game{},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Get Game"), func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/game/%s", tc.Name), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.ExpectedStatus, w.Code)

			if w.Code == http.StatusOK {
				var response *models.Game
				json.Unmarshal([]byte(w.Body.String()), &response)
				assert.Equal(t, tc.ExpectedGame, response)
			}
		})
	}
}

func TestRemoveGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames []string
	}{
		{
			"quiblyv2",
			http.StatusOK,
			[]string{"quibly"},
		},
		{
			"quiblyv2",
			http.StatusNotFound,
			[]string{},
		},
		{
			"quiblyv3",
			http.StatusNotFound,
			[]string{},
		},
		{
			"another_one",
			http.StatusNotFound,
			[]string{},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Remove Game"), func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/game/%s", tc.Name), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.ExpectedStatus, w.Code)

			if w.Code == http.StatusOK {
				req, _ := http.NewRequest("GET", "/game", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var response []string
				json.Unmarshal([]byte(w.Body.String()), &response)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, tc.ExpectedGameNames, response)
			}
		})
	}
}

func TestDisableGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames models.Game
	}{
		{
			"quibly",
			http.StatusOK,
			models.Game{
				Name:      "quibly",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
				Questions: []models.Question{},
				Enabled:   false,
			},
		},
		{
			"quibly",
			http.StatusConflict,
			models.Game{},
		},
		{
			"quiblyv2",
			http.StatusNotFound,
			models.Game{},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Disable A Game"), func(t *testing.T) {
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/game/%s/disable", tc.Name), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.ExpectedStatus, w.Code)

			if w.Code == http.StatusOK {
				req, _ := http.NewRequest("GET", fmt.Sprintf("/game/%s", tc.Name), nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var response *models.Game
				json.Unmarshal([]byte(w.Body.String()), &response)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, tc.ExpectedGameNames, response)
			}
		})
	}
}

func TestEnableGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames models.Game
	}{
		{
			"quibly",
			http.StatusOK,
			models.Game{
				Name:      "quibly",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/quibly",
				Questions: []models.Question{},
				Enabled:   true,
			},
		},
		{
			"quibly",
			http.StatusConflict,
			models.Game{},
		},
		{
			"quiblyv2",
			http.StatusNotFound,
			models.Game{},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Enable A Game"), func(t *testing.T) {
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/game/%s/enable", tc.Name), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.ExpectedStatus, w.Code)

			if w.Code == http.StatusOK {
				req, _ := http.NewRequest("GET", fmt.Sprintf("/game/%s", tc.Name), nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var response *models.Game
				json.Unmarshal([]byte(w.Body.String()), &response)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, tc.ExpectedGameNames, response)
			}
		})
	}
}
