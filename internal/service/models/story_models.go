package models

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// Story struct to contain information about a user story
type Story struct {
	GameName string          `bson:"game_name"          json:"game_name"`
	Username string          `bson:"username"`
	Question string          `bson:"question"`
	Round    string          `bson:"round,omitempty"`
	Nickname string          `bson:"nickname,omitempty"`
	Answers  StoryAnswerType `bson:"answers"`
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
		Username string
		Question string
		Round    string
		Nickname string
	}{}

	var answers struct {
		Answers bson.Raw
	}

	err := unmarshalBSONToStruct(data, &temp, &answers)
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
		Username string
		Question string
		Round    string
		Nickname string
	}{}

	var answers struct {
		Answers json.RawMessage
	}

	err := unmarshalJSONToStruct(data, &temp, &answers)
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
	Username string
	Question string
	Round    string
	Nickname string
}, story *Story) {
	story.GameName = temp.GameName
	story.Username = temp.Username
	story.Question = temp.Question
	story.Round = temp.Round
	story.Nickname = temp.Nickname
}

func getStoryType(gameName string) (StoryAnswerType, error) {
	switch gameName {
	case DRAWLOSSEUM:
		return &StoryDrawlosseumAnswers{}, nil
	case QUIBLY:
		return &StoryQuiblyAnswers{}, nil
	case FIBBINGIT:
		return &StoryFibbingItAnswers{}, nil
	default:
		return nil, errors.Errorf("unknown game name %s", gameName)
	}
}

// StoryFibbingIt contains information about the Fibbing It answers for user stories
type StoryFibbingIt struct {
	Nickname string `bson:"nickname"`
	Answer   string `bson:"answer"`
}

// StoryFibbingItAnswers stores the data for a single fibbing it story.
type StoryFibbingItAnswers []StoryFibbingIt

// NewAnswer creates an empty answer for fibbing it stories.
func (f StoryFibbingItAnswers) NewAnswer() {}

// StoryQuibly contains information about the Quibly answers for user stories
type StoryQuibly struct {
	Nickname string `bson:"nickname"`
	Answer   string `bson:"answer"`
	Votes    int    `bson:"votes"`
}

// StoryQuiblyAnswers stores the data for a single drawlosseum story.
type StoryQuiblyAnswers []StoryQuibly

// NewAnswer creates an empty answer for quibly stories.
func (q StoryQuiblyAnswers) NewAnswer() {}

// StoryDrawlosseum contains information about the Drawlosseum answers for user stories
type StoryDrawlosseum struct {
	Start DrawlosseumDrawingPoint `bson:"start"`
	End   DrawlosseumDrawingPoint `bson:"end"`
	Color string                  `bson:"color"`
}

// DrawlosseumDrawingPoint contains information about a point in a Drawlosseum drawing
type DrawlosseumDrawingPoint struct {
	X float32 `bson:"x"`
	Y float32 `bson:"y"`
}

// StoryDrawlosseumAnswers stores the data for a single drawlosseum story.
type StoryDrawlosseumAnswers []StoryDrawlosseum

// NewAnswer creates an empty answer for drawlosseum stories.
func (d StoryDrawlosseumAnswers) NewAnswer() {}

// Stories is a list of stories.
type Stories []Story

// Add method adds (a list of) stories at once.
func (stories *Stories) Add(db database.Database) error {
	err := db.InsertMultiple("story", stories)
	return err
}

// Get method gets all the stories in the story collection.
func (stories *Stories) Get(db database.Database, filter map[string]string) error {
	err := db.GetAll("story", filter, stories)
	return err
}

// Delete is used to delete a list of stories that match a filter.
func (stories Stories) Delete(db database.Database, filter map[string]string) (bool, error) {
	deleted, err := db.DeleteAll("story", filter)
	return deleted, err
}

// ToInterface converts stories (list of stories) into a list of interfaces, required by GetAll MongoDB.
func (stories Stories) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(stories))
	for i, item := range stories {
		interfaceObject[i] = item
	}
	return interfaceObject
}
