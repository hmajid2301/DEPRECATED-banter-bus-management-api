package models

import (
	"encoding/json"

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

// Add is used to add users to a database.
func (user *User) Add(db Repository) (bool, error) {
	inserted, err := db.Insert("user", user)
	return inserted, err
}

// Get is used to retrieve a user from the database.
func (user *User) Get(db Repository, filter map[string]string) error {
	err := db.Get("user", filter, user)
	return err
}

// Update is used to update a user in the database.
func (user *User) Update(db Repository, filter map[string]string) (bool, error) {
	updated, err := db.Update("user", filter, user)
	return updated, err
}

// Users is a list of users.
type Users []User

// Add method adds (a list of) users at once.
func (users *Users) Add(db Repository) error {
	err := db.InsertMultiple("user", users)
	return err
}

// Get method gets all the users in the user collection.
func (users *Users) Get(db Repository) error {
	err := db.GetAll("user", users)
	return err
}

// ToInterface converts users (list of users) into a list of interfaces, required by GetAll MongoDB.
func (users Users) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(users))
	for i, item := range users {
		interfaceObject[i] = item
	}
	return interfaceObject
}

// Story struct to contain information about a user story
type Story struct {
	GameName string          `bson:"game_name"          json:"game_name"`
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
	Question string
	Round    string
	Nickname string
}, story *Story) {
	story.GameName = temp.GameName
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

// EmptyAnswer creates an empty answer for fibbing it stories.
func (f StoryFibbingItAnswers) EmptyAnswer() {}

// StoryQuibly contains information about the Quibly answers for user stories
type StoryQuibly struct {
	Nickname string `bson:"nickname"`
	Answer   string `bson:"answer"`
	Votes    int    `bson:"votes"`
}

// StoryQuiblyAnswers stores the data for a single drawlosseum story.
type StoryQuiblyAnswers []StoryQuibly

// EmptyAnswer creates an empty answer for quibly stories.
func (q StoryQuiblyAnswers) EmptyAnswer() {}

// StoryDrawlosseum contains information about the Drawlosseum answers for user stories
type StoryDrawlosseum struct {
	Start DrawlosseumDrawingPoint `bson:"start"`
	End   DrawlosseumDrawingPoint `bson:"end"`
	Color string                  `bson:"color"`
}

// StoryDrawlosseumAnswers stores the data for a single drawlosseum story.
type StoryDrawlosseumAnswers []StoryDrawlosseum

// EmptyAnswer creates an empty answer for drawlosseum stories.
func (d StoryDrawlosseumAnswers) EmptyAnswer() {}

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

// EmptyPoolQuestions makes all the fields empty for drawlosseum question pools.
func (d *DrawlosseumQuestionsPool) EmptyPoolQuestions() {
	d.Drawings = []string{}
}

// QuiblyQuestionsPool is the question pool data for quibly.
type QuiblyQuestionsPool struct {
	Pair    []string `bson:"pair,omitempty"`
	Answers []string `bson:"answers,omitempty"`
	Group   []string `bson:"group,omitempty"`
}

// EmptyPoolQuestions makes all the fields empty for fibbing it question pools.
func (f *FibbingItQuestionsPool) EmptyPoolQuestions() {
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

// EmptyPoolQuestions makes all the fields empty for quibly question pools.
func (q *QuiblyQuestionsPool) EmptyPoolQuestions() {
	q.Answers = []string{}
	q.Pair = []string{}
	q.Group = []string{}
}

// NewQuestionPool data required to add a new question pool.
type NewQuestionPool map[string]QuestionPool

// AddToList adds a new question pool to a user.
func (questionPool *NewQuestionPool) AddToList(db Repository, filter map[string]string) (bool, error) {
	updated, err := db.AppendToList("user", filter, questionPool)
	return updated, err
}

// UpdateUserObject is how we can update an existing user object.
type UpdateUserObject map[string]interface{}

// AddToList adds an item (related to a user) to a list.
func (u *UpdateUserObject) AddToList(db Repository, filter map[string]string) (bool, error) {
	updated, err := db.AppendToList("user", filter, u)
	return updated, err
}

// RemoveFromList removes an item (related to a user) from a list.
func (u *UpdateUserObject) RemoveFromList(db Repository, filter map[string]string) (bool, error) {
	updated, err := db.RemoveFromList("user", filter, u)
	return updated, err
}

// QuestionPools data used to store a list of question pools.
type QuestionPools []QuestionPool

// Get method is used to get specific question pools from the database.
func (q *QuestionPools) Get(db Repository, filter map[string]string, parentField string, condition []string) error {
	err := db.GetSubObject("user", filter, parentField, condition, q)
	return err
}
