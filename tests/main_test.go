package controllers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/repository"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/server"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/server/controllers"

	"github.com/gavv/httpexpect"
	"github.com/houqp/gtest"
)

type Tests struct {
	httpExpect *httpexpect.Expect
	DB         core.Repository
}

type TestData struct {
	Games []models.GameInfo `json:"games"`
	Users []models.User     `json:"users"`
}

func (s *Tests) Setup(t *testing.T) {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "config.test.yml")
	config, err := core.NewConfig()
	if err != nil {
		fmt.Printf("Failed to load config %s", err)
	}
	logger := core.SetupLogger(ioutil.Discard)
	core.UpdateLogLevel(logger, "DEBUG")
	db, err := repository.NewMongoDB(logger,
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.DatabaseName,
		config.Database.MaxConns,
		config.Database.Timeout)

	if err != nil {
		fmt.Println(err)
		return
	}

	env := &controllers.Env{Logger: logger, Config: config, DB: db}
	router, err := server.SetupWebServer(env)
	if err != nil {
		fmt.Printf("Failed to setup web server %s", err)
	}

	s.DB = db
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
	InsertData(s.DB, "data/json/game_collection.json", "game")
	InsertData(s.DB, "data/json/user_collection.json", "user")
}

func (s *Tests) AfterEach(t *testing.T) {
	err := s.DB.RemoveCollection("game")
	if err != nil {
		fmt.Printf("Failed to remove collection game %s", err)
	}

	err = s.DB.RemoveCollection("user")
	if err != nil {
		fmt.Printf("Failed to remove collection user %s", err)
	}
}

func TestSampleTests(t *testing.T) {
	gtest.RunSubTests(t, &Tests{})
}

func InsertData(db core.Repository, dataFilePath string, collection string) {
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

	err = db.InsertMultiple(collection, dataList)
	if err != nil {
		fmt.Println(err)
	}
}
