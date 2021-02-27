package models

import "gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"

// Question is the data for questions related to a game.
type Question struct {
	Content map[string]string `bson:"content"`
	Enabled *bool             `bson:"enabled,omitempty"`
}

// NewQuestion type is used to add new questions to games.
type NewQuestion map[string]Question

// AddToList adds a new Question to the game.
func (question *NewQuestion) AddToList(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.AppendToList("game", filter, question)
	return updated, err
}

// QuiblyQuestions is the data for questions related to the Quibly game.
type QuiblyQuestions struct {
	Pair    []Question `bson:"pair,omitempty"`
	Answers []Question `bson:"answers,omitempty"`
	Group   []Question `bson:"group,omitempty"`
}

// NewQuestions creates an empty quibly it questions.
func (q *QuiblyQuestions) NewQuestions() {
	q.Answers = []Question{}
	q.Pair = []Question{}
	q.Group = []Question{}
}

// DrawlosseumQuestions is the data required to play the Drawlosseum game.
type DrawlosseumQuestions struct {
	Drawings []Question `bson:"drawings"`
}

// NewQuestions creates an empty drawlosseum it questions.
func (d *DrawlosseumQuestions) NewQuestions() {
	d.Drawings = []Question{}
}

// FibbingItQuestions is the data for questions related to the Fibbing It game.
// NOTE: Some fields are a map of questions, because questions will be grouped with other similar questions
type FibbingItQuestions struct {
	Opinion  map[string]map[string][]Question `bson:"opinion"`
	FreeForm map[string][]Question            `bson:"free_form" json:"free_form"`
	Likely   []Question                       `bson:"likely"`
}

// NewQuestions creates an empty fibbing it questions.
func (f *FibbingItQuestions) NewQuestions() {
	f.Opinion = map[string]map[string][]Question{}
	f.FreeForm = map[string][]Question{}
	f.Likely = []Question{}
}

// GenericQuestion is generic structure all questions can take, has all the required fields for any question.
type GenericQuestion struct {
	Content      string
	Round        string
	LanguageCode string
	Group        *GenericQuestionGroup
}

// GenericQuestionGroup provides extra context to a question, when it belong to a group.
type GenericQuestionGroup struct {
	Name string
	Type string
}
