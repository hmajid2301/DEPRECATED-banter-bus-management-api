package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/houqp/gtest"
	"github.com/wI2L/fizz"
	"gopkg.in/go-playground/assert.v1"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server"
	"banter-bus-server/src/server/models"
	"banter-bus-server/src/utils/config"
)

var router *fizz.Fizz

type SampleTests struct{}

type TestData struct {
	Games []models.Game `json:"games"`
}

func (s *SampleTests) Setup(t *testing.T) {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "config.test.yml")
	os.Setenv("BANTER_BUS_LOG_LEVEL", "FATAL")
	config := config.GetConfig()
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
func (s *SampleTests) Teardown(t *testing.T) {}

func (s *SampleTests) BeforeEach(t *testing.T) {
	data, _ := ioutil.ReadFile("test_data.json")
	var docs TestData
	json.Unmarshal(data, &docs)
	var ui []interface{}
	for _, t := range docs.Games {
		ui = append(ui, t)
	}

	database.InsertMultiple("game", ui)
}

func (s *SampleTests) AfterEach(t *testing.T) {
	database.RemoveCollection("game")
}

func TestSampleTests(t *testing.T) {
	gtest.RunSubTests(t, &SampleTests{})
}

func (s *SampleTests) SubTestAddGame(t *testing.T) {
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

func (s *SampleTests) SubTestGetAllGames(t *testing.T) {
	req, _ := http.NewRequest("GET", "/game", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var response []string
	json.Unmarshal([]byte(w.Body.String()), &response)

	var expectedResult = []string{
		"a_game",
		"fibbly",
		"draw_me",
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResult, response)
}

func (s *SampleTests) SubTestGetGame(t *testing.T) {
	cases := []struct {
		Name           string
		ExpectedStatus int
		ExpectedGame   models.Game
	}{
		{
			"a_game",
			http.StatusOK,
			models.Game{
				Name:      "a_game",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/a_game",
				Questions: []models.Question{},
				Enabled:   true,
			},
		},
		{
			"fibbly",
			http.StatusOK,
			models.Game{
				Name:      "fibbly",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/fibbly",
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

func (s *SampleTests) SubTestRemoveGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames []string
	}{
		{
			"a_game",
			http.StatusOK,
			[]string{"fibbly", "draw_me"},
		},
		{
			"a_game",
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

func (s *SampleTests) SubTestEnableGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames models.Game
	}{
		{
			"draw_me",
			http.StatusOK,
			models.Game{
				Name:      "draw_me",
				RulesURL:  "https://google.com/draw_me",
				Questions: []models.Question{},
				Enabled:   true,
			},
		},
		{
			"draw_me",
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

func (s *SampleTests) SubTestDisableGame(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedGameNames models.Game
	}{
		{
			"a_game",
			http.StatusOK,
			models.Game{
				Name:      "a_game",
				RulesURL:  "https://gitlab.com/banter-bus/banter-bus-server/-/wikis/docs/rules/a_game",
				Questions: []models.Question{},
				Enabled:   false,
			},
		},
		{
			"a_game",
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

func (s *SampleTests) SubTestGetOpenAPI(t *testing.T) {
	req, _ := http.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var a = []byte(w.Body.String())
	var out bytes.Buffer
	json.Indent(&out, a, "", "  ")
	ioutil.WriteFile("../openapi.json", out.Bytes(), 0644)
}
