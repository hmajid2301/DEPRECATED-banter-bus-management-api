package models

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/mgo.v2/bson"
)

// User struct to hold a database entry for a user, i.e a player with an account
type User struct {
	Username      string           `bson:"username"`
	Admin         *bool            `bson:"admin,omitempty"`
	Privacy       string           `bson:"privacy,omitempty"`
	Membership    string           `bson:"membership,omitempty"`
	Preferences   *UserPreferences `bson:"preferences,omitempty"`
	Friends       []Friend         `bson:"friends"`
	Stories       []Story          `bson:"stories"`
	QuestionPools []QuestionPool   `bson:"question_pools"        json:"question_pools"`
}

// Story struct to contain information about a user story
type Story struct {
	GameName string      `bson:"game_name"          json:"game_name"`
	Question string      `bson:"question"`
	Round    string      `bson:"round,omitempty"`
	Nickname string      `bson:"nickname,omitempty"`
	Answers  interface{} `bson:"answers"`
}

// StoryFibbingIt contains information about the Fibbing It answers for user stories
type StoryFibbingIt struct {
	Nickname string `bson:"nickname"`
	Answer   string `bson:"answer"`
}

// StoryQuibly contains information about the Quibly answers for user stories
type StoryQuibly struct {
	Nickname string `bson:"nickname"`
	Answer   string `bson:"answer"`
	Votes    int    `bson:"votes"`
}

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

// UserPreferences struct to hold information on a user's preferences
type UserPreferences struct {
	LanguageCode string `bson:"language_code,omitempty"`
}

// Friend struct to hold details about a user's friend
type Friend struct {
	Username string `bson:"username"`
}

// QuestionPool struct holds the data for all their own questions
type QuestionPool struct {
	PoolName  string      `bson:"pool_name" json:"pool_name"`
	GameName  string      `bson:"game_name" json:"game_name"`
	Privacy   string      `bson:"privacy"`
	Questions interface{} `bson:"questions"`
}

// DrawlosseumQuestionsPool is the question pool data for drawlosseum.
type DrawlosseumQuestionsPool struct {
	Drawings []string `bson:"drawings,omitempty"`
}

// QuiblyQuestionsPool is the question pool data for quibly.
type QuiblyQuestionsPool struct {
	Pair    []string `bson:"pair,omitempty"`
	Answers []string `bson:"answers,omitempty"`
	Group   []string `bson:"group,omitempty"`
}

// FibbingItQuestionsPool is the question pool data for quibly.
type FibbingItQuestionsPool struct {
	Opinion  map[string]map[string][]string `bson:"opinion,omitempty"`
	FreeForm map[string][]string            `bson:"free_form,omitempty"`
	Likely   []string                       `bson:"likely,omitempty"`
}

// UnmarshalBSONValue is a custom unmarshall function, that will unmarshall question pools differently, i.e. questions
// field depending on the game name. As each has it's own question structure. The main purpose is just to get
// the raw BSON data for the `Questions` field.
//
// Then, we work out the type of game, using `GameName`. Depending on the game we unmarshall the data into different
// structs and assign that to the `Questions` field of that `questionPool` variable, which is of type `QuestionPool`.
// The `questionPool`, is what is returned when we get the `QuestionPool` data from the database.
func (questionPool *QuestionPool) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var questions struct {
		Questions bson.Raw
	}

	err := unmarshalBSONToStruct(data, &questionPool, &questions)
	if err != nil {
		return err
	}

	switch questionPool.GameName {
	case "drawlosseum":
		questionStructure := DrawlosseumQuestionsPool{}
		err = questions.Questions.Unmarshal(&questionStructure)
		questionPool.Questions = questionStructure
	case "quibly":
		questionStructure := QuiblyQuestionsPool{}
		err = questions.Questions.Unmarshal(&questionStructure)
		questionPool.Questions = questionStructure
	case "fibbing_it":
		questionStructure := FibbingItQuestionsPool{}
		err = questions.Questions.Unmarshal(&questionStructure)
		questionPool.Questions = questionStructure
	default:
		return errors.Errorf("Unknown game name %s", questionPool.GameName)
	}

	return err
}

// UnmarshalBSONValue is a custom unmarshall function, that will unmarshall question pools differently, i.e. answers
// field depending on the game name. As each has it's own question structure. The main purpose is just to get
// the raw BSON data for the `Answers` field.
//
// Then, we work out the type of game, using `GameName`. Depending on the game we unmarshall the data into different
// structs and assign that to the `Questions` field of that `story` variable, which is of type `Story`.
// The `story`, is what is returned when we get the `Story` data from the database.
func (story *Story) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var answers struct {
		Answers bson.Raw
	}

	err := unmarshalBSONToStruct(data, &story, &answers)
	if err != nil {
		return err
	}

	switch story.GameName {
	case "drawlosseum":
		storyStructure := []StoryDrawlosseum{}
		err = answers.Answers.Unmarshal(&storyStructure)
		story.Answers = storyStructure
	case "quibly":
		storyStructure := []StoryQuibly{}
		err = answers.Answers.Unmarshal(&storyStructure)
		story.Answers = storyStructure
	case "fibbing_it":
		storyStructure := []StoryFibbingIt{}
		err = answers.Answers.Unmarshal(&storyStructure)
		story.Answers = storyStructure
	default:
		return errors.Errorf("Unknown game name %s", story.GameName)
	}

	return err
}

// unmarshalBSONToStruct is a custom unmarshal function, that will unmarshal BSON to structs.
// This function is to be used when a subfield can take multiple types i.e. storing question
// data differently for different games. It will unmarshal that field (polymorphic one)
// i.e. `Questions`, into BSON raw data. This can then be cast into the correct struct for
// the polymorphic field.
//
// The first unmarshal gets the Raw BSON data. The Raw BSON data allows us to unmarshal sub-objects like `Questions`
// field to a specific struct.
//
// The second unmarshal converts the raw BSON data into a struct i.e. `QuestionPool`, note in this example `Questions`
//  field is type `interface{}`.
//
// Next, we unmarshal the subField into raw BSON data, in the example above this would be the `Questions` field.
// This way we only have raw BSON data related to that field and can be cast appropriate.
func unmarshalBSONToStruct(data []byte, structType interface{}, subField interface{}) error {
	var rawData bson.Raw
	err := bson.Unmarshal(data, &rawData)
	if err != nil {
		return err
	}

	err = rawData.Unmarshal(structType)
	if err != nil {
		return err
	}

	err = rawData.Unmarshal(subField)
	if err != nil {
		return err
	}

	return nil
}
