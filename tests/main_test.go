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
	DB         models.Repository
}

type GameTestData struct {
	Games models.Games `json:"games"`
}

type UserTestData struct {
	Users models.Users `json:"users"`
}

type TestData interface {
	InsertData(path string, db models.Repository)
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
	gameData := GameTestData{}
	gameData.InsertData("data/json/game_collection.json", s.DB)

	userData := UserTestData{}
	userData.InsertData("data/json/user_collection.json", s.DB)
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

func (gameData *GameTestData) InsertData(path string, db models.Repository) {
	err := getData(path, gameData)
	if err != nil {
		fmt.Println(err)
	}

	err = gameData.Games.Add(db)
	if err != nil {
		fmt.Println(err)
	}
}

func (userData *UserTestData) InsertData(path string, db models.Repository) {
	err := getData(path, userData)
	if err != nil {
		fmt.Println(err)
	}

	err = userData.Users.Add(db)
	if err != nil {
		fmt.Println(err)
	}
}

func getData(path string, model TestData) error {
	data, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(data, &model)
	return err
}
