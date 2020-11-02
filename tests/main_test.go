package controllers_test

import (
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server"
	"banter-bus-server/src/server/models"
	"banter-bus-server/src/utils/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/houqp/gtest"
	"github.com/wI2L/fizz"
)

type Tests struct {
	router *fizz.Fizz
}

type TestData struct {
	Games []models.Game `json:"games"`
}

func (s *Tests) Setup(t *testing.T) {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "config.test.yml")
	log.SetOutput(ioutil.Discard)
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
		docs     TestData
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
