package models

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// QuestionPool struct holds the data for all their own questions
type QuestionPool struct {
	PoolName     string           `bson:"pool_name"     json:"pool_name"`
	GameName     string           `bson:"game_name"     json:"game_name"`
	LanguageCode string           `bson:"language_code" json:"language_code"`
	Privacy      string           `bson:"privacy"`
	Questions    QuestionPoolType `bson:"questions"`
}

// UnmarshalBSONValue is a custom unmarshal function, that will unmarshal question pools differently, i.e. questions
// field depending on the game name. As each has it's own question structure. The main purpose is just to get
// the raw BSON data for the `Questions` field.
//
// Then, we work out the type of game, using `GameName`. Depending on the game we unmarshal the data into different
// structs and assign that to the `Questions` field of that `questionPool` variable, which is of type `QuestionPool`.
// The `questionPool`, is what is returned when we get the `QuestionPool` data from the database.
func (questionPool *QuestionPool) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	temp := struct {
		PoolName     string `bson:"pool_name" json:"pool_name"`
		GameName     string `bson:"game_name" json:"game_name"`
		LanguageCode string `bson:"language_code" json:"language_code"`
		Privacy      string
	}{}

	var questions struct {
		Questions bson.Raw
	}

	err := unmarshalBSONToStruct(data, &temp, &questions)
	if err != nil {
		return err
	}

	setQuestionPoolFields(temp, questionPool)
	questionStructure, err := getQuestionPoolType(questionPool.GameName)
	if err != nil {
		return err
	}

	err = questions.Questions.Unmarshal(questionStructure)
	if err != nil {
		return err
	}

	questionPool.Questions = questionStructure
	return nil
}

// UnmarshalJSON works the same way as the UnmarshalBSONToStruct function above.
func (questionPool *QuestionPool) UnmarshalJSON(data []byte) error {
	temp := struct {
		PoolName     string `bson:"pool_name" json:"pool_name"`
		GameName     string `bson:"game_name" json:"game_name"`
		LanguageCode string `bson:"language_code" json:"language_code"`
		Privacy      string
	}{}

	var questions struct {
		Questions json.RawMessage
	}

	err := unmarshalJSONToStruct(data, &temp, &questions)
	if err != nil {
		return err
	}

	setQuestionPoolFields(temp, questionPool)
	questionStructure, err := getQuestionPoolType(questionPool.GameName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(questions.Questions, &questionStructure)
	if err != nil {
		return err
	}

	questionPool.Questions = questionStructure
	return nil
}

func setQuestionPoolFields(temp struct {
	PoolName     string `bson:"pool_name" json:"pool_name"`
	GameName     string `bson:"game_name" json:"game_name"`
	LanguageCode string `bson:"language_code" json:"language_code"`
	Privacy      string
}, questionPool *QuestionPool) {
	questionPool.PoolName = temp.PoolName
	questionPool.GameName = temp.GameName
	questionPool.LanguageCode = temp.LanguageCode
	questionPool.Privacy = temp.Privacy
}

func getQuestionPoolType(gameName string) (QuestionPoolType, error) {
	switch gameName {
	case DRAWLOSSEUM:
		return &DrawlosseumQuestionsPool{}, nil
	case QUIBLY:
		return &QuiblyQuestionsPool{}, nil
	case FIBBINGIT:
		return &FibbingItQuestionsPool{}, nil
	default:
		return nil, errors.Errorf("unknown game name %s", gameName)
	}
}

// DrawlosseumQuestionsPool is the question pool data for drawlosseum.
type DrawlosseumQuestionsPool struct {
	Drawings []string `bson:"drawings,omitempty"`
}

// NewPool makes all the fields empty for drawlosseum question pools.
func (d *DrawlosseumQuestionsPool) NewPool() {
	d.Drawings = []string{}
}

// QuiblyQuestionsPool is the question pool data for quibly.
type QuiblyQuestionsPool struct {
	Pair    []string `bson:"pair,omitempty"`
	Answers []string `bson:"answers,omitempty"`
	Group   []string `bson:"group,omitempty"`
}

// NewPool makes all the fields empty for fibbing it question pools.
func (f *FibbingItQuestionsPool) NewPool() {
	f.Opinion = map[string]map[string][]string{}
	f.FreeForm = map[string][]string{}
	f.Likely = []string{}
}

// FibbingItQuestionsPool is the question pool data for quibly.
type FibbingItQuestionsPool struct {
	Opinion  map[string]map[string][]string `bson:"opinion,omitempty"`
	FreeForm map[string][]string            `bson:"free_form,omitempty" json:"free_form,omitempty"`
	Likely   []string                       `bson:"likely,omitempty"`
}

// NewPool makes all the fields empty for quibly question pools.
func (q *QuiblyQuestionsPool) NewPool() {
	q.Answers = []string{}
	q.Pair = []string{}
	q.Group = []string{}
}

// NewPoolQuestion data required to add a new question pool.
type NewPoolQuestion map[string]QuestionPool

// AddToList adds a new question pool to a user.
func (q *NewPoolQuestion) AddToList(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.AppendToList("user", filter, q)
	return updated, err
}

// UpdateUserObject is how we can update an existing user object.
type UpdateUserObject map[string]interface{}

// AddToList adds an item (related to a user) to a list.
func (u *UpdateUserObject) AddToList(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.AppendToList("user", filter, u)
	return updated, err
}

// RemoveFromList removes an item (related to a user) from a list.
func (u *UpdateUserObject) RemoveFromList(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.RemoveFromList("user", filter, u)
	return updated, err
}

// QuestionPools data used to store a list of question pools.
type QuestionPools []QuestionPool

// Get method is used to get specific question pools from the database.
func (q *QuestionPools) Get(
	db database.Database,
	filter map[string]string,
	parentField string,
	condition []string,
) error {
	err := db.GetSubObject("user", filter, parentField, condition, q)
	return err
}
