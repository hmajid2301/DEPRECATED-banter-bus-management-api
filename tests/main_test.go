package controllers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/models"
	"banter-bus-server/src/server"
	"banter-bus-server/src/utils/config"

	"github.com/gavv/httpexpect"
	"github.com/houqp/gtest"
	"github.com/sirupsen/logrus"
)

type Tests struct {
	httpExpect *httpexpect.Expect
}

type TestData struct {
	Games []models.Game `json:"games"`
	Users []models.User `json:"users"`
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
	router, err := server.NewRouter()

	if err != nil {
		fmt.Println(err)
		return
	}
	s.httpExpect = httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(router.Engine()),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func (s *Tests) Teardown(t *testing.T) {}

func (s *Tests) BeforeEach(t *testing.T) {
	InsertData("data/json/game_collection.json", "game")
	InsertData("data/json/user_collection.json", "user")
}

func (s *Tests) AfterEach(t *testing.T) {
	database.RemoveCollection("game")
	database.RemoveCollection("user")
}

func TestSampleTests(t *testing.T) {
	gtest.RunSubTests(t, &Tests{})
}

func InsertData(dataFilePath string, collection string) {
	data, _ := ioutil.ReadFile(dataFilePath)

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

	for _, t := range docs.Users {
		dataList = append(dataList, t)
	}

	err = database.InsertMultiple(collection, dataList)
	if err != nil {
		fmt.Println(err)
	}
}
