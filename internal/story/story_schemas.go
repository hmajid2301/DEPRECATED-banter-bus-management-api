package story

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/banter-bus/banter-bus-management-api/internal"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/games"
)

// Story struct to contain information about a user story
type Story struct {
	GameName string          `bson:"game_name"          json:"game_name"`
	ID       string          `bson:"id"`
	Question string          `bson:"question"`
	Round    string          `bson:"round,omitempty"`
	Nickname string          `bson:"nickname,omitempty"`
	Answers  StoryAnswerType `bson:"answers"`
}

func (story *Story) Add(db database.Database) (bool, error) {
	inserted, err := db.Insert("story", story)
	return inserted, err
}

func (story *Story) Get(db database.Database, filter map[string]interface{}) error {
	err := db.Get("story", filter, story)
	return err
}

func (story *Story) Update(db database.Database, filter map[string]interface{}) (bool, error) {
	updated, err := db.Update("story", filter, story)
	return updated, err
}

// UnmarshalBSONValue is a custom unmarshal function, that will unmarshal question pools differently, i.e. answers
// field depending on the game name. As each has it's own question structure. The main purpose is just to get
// the raw BSON data for the `Answers` field.
//
// Then, we work out the type of game, using `GameName`. Depending on the game we unmarshal the data into different
// structs and assign that to the `Questions` field of that `story` variable, which is of type `Story`.
// The `story`, is what is returned when we get the `Story` data from the database.
func (story *Story) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	temp := struct {
		GameName string `json:"game_name" bson:"game_name"`
		ID       string
		Question string
		Round    string
		Nickname string
	}{}

	var answers struct {
		Answers bson.Raw
	}

	err := internal.UnmarshalBSONToStruct(data, &temp, &answers)
	if err != nil {
		return err
	}

	setStoryFields(temp, story)
	storyStructure, err := getStoryType(temp.GameName)
	if err != nil {
		return err
	}

	err = answers.Answers.Unmarshal(storyStructure)
	story.Answers = storyStructure
	return err
}

// UnmarshalJSON works almost the same way as the UnmarshalBSONValue method above.
func (story *Story) UnmarshalJSON(data []byte) error {
	temp := struct {
		GameName string `json:"game_name" bson:"game_name"`
		ID       string
		Question string
		Round    string
		Nickname string
	}{}

	var answers struct {
		Answers json.RawMessage
	}

	err := internal.UnmarshalJSONToStruct(data, &temp, &answers)
	if err != nil {
		return err
	}

	setStoryFields(temp, story)
	storyStructure, err := getStoryType(story.GameName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(answers.Answers, &storyStructure)
	story.Answers = storyStructure
	return err
}

func setStoryFields(temp struct {
	GameName string `json:"game_name" bson:"game_name"`
	ID       string
	Question string
	Round    string
	Nickname string
}, story *Story) {
	story.GameName = temp.GameName
	story.ID = temp.ID
	story.Question = temp.Question
	story.Round = temp.Round
	story.Nickname = temp.Nickname
}

func getStoryType(gameName string) (StoryAnswerType, error) {
	switch gameName {
	case games.DRAWLOSSEUM:
		return &DrawlosseumAnswers{}, nil
	case games.QUIBLY:
		return &QuiblyAnswers{}, nil
	case games.FIBBINGIT:
		return &FibbingItAnswers{}, nil
	default:
		return nil, errors.Errorf("unknown game name %s", gameName)
	}
}

type FibbingItAnswer struct {
	Nickname string `bson:"nickname"`
	Answer   string `bson:"answer"`
}

type FibbingItAnswers []FibbingItAnswer

func (f FibbingItAnswers) NewAnswer() {}

type QuiblyAnswer struct {
	Nickname string `bson:"nickname"`
	Answer   string `bson:"answer"`
	Votes    int    `bson:"votes"`
}

type QuiblyAnswers []QuiblyAnswer

func (q QuiblyAnswers) NewAnswer() {}

type DrawlosseumAnswers []CaertsianCoordinateColor

func (d DrawlosseumAnswers) NewAnswer() {}

type Stories []Story

func (stories *Stories) Add(db database.Database) error {
	err := db.InsertMultiple("story", stories)
	return err
}

func (stories *Stories) Get(db database.Database, filter map[string]interface{}) error {
	err := db.GetAll("story", filter, stories)
	return err
}

func (stories *Stories) GetWithLimit(db database.Database, filter map[string]interface{}, limit int64) error {
	return nil
}

func (stories Stories) Delete(db database.Database, filter map[string]interface{}) (bool, error) {
	deleted, err := db.DeleteAll("story", filter)
	return deleted, err
}

func (stories Stories) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(stories))
	for i, item := range stories {
		interfaceObject[i] = item
	}
	return interfaceObject
}

type StoryAnswerType interface {
	NewAnswer()
}
