package controllers_test

import (
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/dbmodels"
	"banter-bus-server/src/server"
	"banter-bus-server/src/utils/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/houqp/gtest"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
)

type Tests struct {
	router *fizz.Fizz
}

type GameData struct {
	Name      string             `bson:"name"`
	Questions *dbmodels.Question `bson:"questions"`
	RulesURL  string             `json:"rules_url" bson:"rules_url"`
	Enabled   bool               `bson:"enabled"`
}

type TestData struct {
	Games []GameData `json:"games"`
}

func (s *Tests) Setup(t *testing.T) {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "config.test.yml")
	logrus.SetOutput(ioutil.Discard)
	config := config.GetConfig()
	dbConfig := database.Config{
		Username:     config.Database.Username,
		Password:     config.Database.Password,
		DatabaseName: config.Database.DatabaseName,
		Host:         config.Database.Host,
		Port:         config.Database.Port,
	}
	database.InitialiseDatabase(dbConfig)
	router, _ := server.NewRouter()
	s.router = router
}

func (s *Tests) Teardown(t *testing.T) {}

func (s *Tests) BeforeEach(t *testing.T) {
	InsertData("data/game.json", "game")
}

func (s *Tests) AfterEach(t *testing.T) {
	database.RemoveCollection("game")
}

func TestSampleTests(t *testing.T) {
	gtest.RunSubTests(t, &Tests{})
}

func InsertData(dataFilePath string, collection string) {
	data, _ := ioutil.ReadFile("data/game.json")
	var (
		docs     *TestData
		dataList []interface{}
	)

	err := json.Unmarshal(data, &docs)
	if err != nil {
		fmt.Println(err)
	}

	for _, t := range docs.Games {
		dataList = append(dataList, t)
	}

	err = database.InsertMultiple(collection, dataList)
	if err != nil {
		fmt.Println(err)
	}
}
