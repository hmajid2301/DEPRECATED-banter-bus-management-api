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
	Friends       []Friend         `bson:"friends,omitempty"`
	Stories       []Story          `bson:"stories,omitempty"`
	QuestionPools []QuestionPool   `bson:"question_pools,omitempty" json:"question_pools,omitempty"`
}

// Story struct to contain information about a user story
type Story struct {
	GameName string        `bson:"game_name" json:"game_name"`
	Question string        `bson:"question"`
	Answers  []StoryAnswer `bson:"answers"`
}

// StoryAnswer struct to contain information needed to store an answer for a player's story
type StoryAnswer struct {
	Answer string `bson:"answer"`
	Votes  int    `bson:"votes"`
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
// field depending on the game name. As each has it's own question structure.
//
// The first unmarshall gets the Raw BSON data. The Raw BSON data allows us to unmarshall sub-objects like `Questions`
// field to a specific struct.
//
// The second unmarshall converts the raw BSON data into a `QuestionPool` struct, note the `Questions` field is type
// `interface{}`. However we need it to be one of the types of the game.
//
// Next, we unmarshall the `Questions` field into raw BSON data. This way we only have raw BSON data related to
// `Questions` field.
//
// Finally, we work out the type of game, using `GameName`. Depending on the game we unmarshall the data into different
// structs and assign that to the `Questions` field of that `questionPool` variable, which is of type `QuestionPool`.
// The `questionPool`, is what is returned when we get the `QuestionPool` data from the database.
func (questionPool *QuestionPool) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var rawData bson.Raw
	err := bson.Unmarshal(data, &rawData)
	if err != nil {
		return err
	}

	err = rawData.Unmarshal(&questionPool)
	if err != nil {
		return err
	}

	var questions struct {
		Questions bson.Raw
	}

	err = rawData.Unmarshal(&questions)
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
