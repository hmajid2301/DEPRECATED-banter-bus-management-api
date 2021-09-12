package controllers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/api"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/questions"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/story"

	"github.com/gavv/httpexpect"
	"github.com/houqp/gtest"
)

type Tests struct {
	httpExpect *httpexpect.Expect
	DB         database.Database
}

type GameTestData struct {
	Games games.Games `json:"games"`
	DB    database.Database
}

type QuestionTestData struct {
	Questions questions.Questions `json:"questions"`
	DB        database.Database
}

type StoryTestData struct {
	Stories story.Stories `json:"stories"`
	DB      database.Database
}

type TestData interface {
	InsertData(path string)
}

func (s *Tests) Setup(t *testing.T) {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "config.test.yml")
	conf, err := core.NewConfig()
	if err != nil {
		fmt.Printf("Failed to load config %s", err)
		return
	}
	logger := core.SetupLogger(ioutil.Discard)
	core.UpdateLogLevel(logger, "DEBUG")
	db, err := database.NewMongoDB(logger,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Name,
		conf.DB.MaxConns,
		conf.DB.Timeout,
		conf.DB.AuthDB)

	if err != nil {
		fmt.Println(err)
		return
	}

	env := &api.Env{Logger: logger, Conf: conf, DB: db}
	router, err := api.Setup(env)
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
	gameData := GameTestData{
		DB: s.DB,
	}
	gameData.InsertData("data/json/game_collection.json")

	questionData := QuestionTestData{
		DB: s.DB,
	}
	questionData.InsertData("data/json/question_collection.json")

	storyData := StoryTestData{
		DB: s.DB,
	}
	storyData.InsertData("data/json/story_collection.json")
}

func (s *Tests) AfterEach(t *testing.T) {
	err := s.DB.RemoveCollection("game")
	if err != nil {
		fmt.Printf("Failed to remove collection game %s", err)
	}

	err = s.DB.RemoveCollection("question")
	if err != nil {
		fmt.Printf("Failed to remove collection question %s", err)
	}

	err = s.DB.RemoveCollection("story")
	if err != nil {
		fmt.Printf("Failed to remove collection story %s", err)
	}
}

func TestSampleTests(t *testing.T) {
	gtest.RunSubTests(t, &Tests{})
}

func (g *GameTestData) InsertData(path string) {
	err := getData(path, g)
	if err != nil {
		fmt.Println(err)
	}

	err = g.Games.Add(g.DB)
	if err != nil {
		fmt.Println(err)
	}
}

func (q *QuestionTestData) InsertData(path string) {
	err := getData(path, q)
	if err != nil {
		fmt.Println(err)
	}

	err = q.Questions.Add(q.DB)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *StoryTestData) InsertData(path string) {
	err := getData(path, s)
	if err != nil {
		fmt.Println(err)
	}

	err = s.Stories.Add(s.DB)
	if err != nil {
		fmt.Println(err)
	}
}

func getData(path string, model TestData) error {
	data, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(data, &model)
	return err
}
