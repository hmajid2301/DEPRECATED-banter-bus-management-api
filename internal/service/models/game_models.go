package models

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// QUIBLY is the string representation of the game name.
const QUIBLY = "quibly"

// FIBBINGIT is the string representation of the game name.
const FIBBINGIT = "fibbing_it"

// DRAWLOSSEUM is the string representation of the game name.
const DRAWLOSSEUM = "drawlosseum"

// Game is the data required by all games. The only thing that varies between different games
// is how they store questions.
type Game struct {
	Name      string       `bson:"name,omitempty"`
	RulesURL  string       `bson:"rules_url,omitempty" json:"rules_url,omitempty"`
	Enabled   *bool        `bson:"enabled,omitempty"`
	Questions QuestionType `bson:"questions,omitempty"`
}

// Add is used to add a game.
func (game *Game) Add(db database.Database) (bool, error) {
	inserted, err := db.Insert("game", game)
	return inserted, err
}

// Get is used to get a game.
func (game *Game) Get(db database.Database, filter map[string]string) error {
	err := db.Get("game", filter, game)
	return err
}

// Update is used to update a game.
func (game *Game) Update(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.Update("game", filter, game)
	return updated, err
}

// HasGroups checks if the game has question groups for the specified round
// More efficient way of storing strings for lookup than a slice hence it's a map
func (game *Game) HasGroups(round string) bool {
	var gameRoundsWithGroups = map[string]struct{}{
		"fibbing_it.opinion":   {},
		"fibbing_it.free_form": {},
	}

	queryString := fmt.Sprintf("%s.%s", game.Name, round)
	_, isPresent := gameRoundsWithGroups[queryString]
	return isPresent
}

// UnmarshalBSONValue is a custom unmarshal function, that will unmarshal question pools differently, i.e. questions
// field depending on the game name. As each has it's own question structure. The main purpose is just to get
// the raw BSON data for the `Questions` field.
//
// Then, we work out the type of game, using `Name`. Depending on the game we unmarshal the data into different
func (game *Game) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	temp := struct {
		Name     string
		RulesURL string `bson:"rules_url" json:"rules_url"`
		Enabled  *bool
	}{}

	var questions struct {
		Questions bson.Raw
	}

	err := unmarshalBSONToStruct(data, &temp, &questions)
	if err != nil {
		return err
	}

	setGameFields(temp, game)
	questionStructure, err := getQuestionType(game.Name)
	if err != nil {
		return err
	}

	err = questions.Questions.Unmarshal(questionStructure)
	if err != nil {
		return err
	}

	game.Questions = questionStructure
	return nil
}

// UnmarshalJSON works almost the same way as the UnmarshalBSONValue method above.
func (game *Game) UnmarshalJSON(data []byte) error {
	temp := struct {
		Name     string
		RulesURL string `bson:"rules_url" json:"rules_url"`
		Enabled  *bool
	}{}

	var questions struct {
		Questions json.RawMessage
	}

	err := unmarshalJSONToStruct(data, &temp, &questions)
	if err != nil {
		return err
	}

	setGameFields(temp, game)
	questionStructure, err := getQuestionType(game.Name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(questions.Questions, &questionStructure)
	if err != nil {
		return err
	}

	game.Questions = questionStructure
	return nil
}

func setGameFields(temp struct {
	Name     string
	RulesURL string `bson:"rules_url" json:"rules_url"`
	Enabled  *bool
}, game *Game) {
	game.Name = temp.Name
	game.RulesURL = temp.RulesURL
	game.Enabled = temp.Enabled
}

func getQuestionType(gameName string) (QuestionType, error) {
	var questionStructure QuestionType

	switch gameName {
	case DRAWLOSSEUM:
		questionStructure = &DrawlosseumQuestions{}
	case QUIBLY:
		questionStructure = &QuiblyQuestions{}
	case FIBBINGIT:
		questionStructure = &FibbingItQuestions{}
	default:
		return nil, errors.Errorf("unknown game name %s", gameName)
	}

	return questionStructure, nil
}

// Games is a list of game(s).
type Games []Game

// Add is used to add a list of games of the database.
func (games *Games) Add(db database.Database) error {
	err := db.InsertMultiple("game", games)
	return err
}

// Get is used to get a list of games.
func (games *Games) Get(db database.Database) error {
	err := db.GetAll("game", games)
	return err
}

// ToInterface converts users (list of users) into a list of interfaces, required by GetAll MongoDB.
func (games Games) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(games))
	for i, item := range games {
		interfaceObject[i] = item
	}
	return interfaceObject
}
