package models

import "gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"

// Question is the data for questions related to a game.
type Question struct {
	GameName string            `bson:"game_name"       json:"game_name"`
	Round    string            `bson:"round,omitempty"`
	Enabled  *bool             `bson:"enabled"`
	Content  map[string]string `bson:"content"`
	Group    QuestionGroup     `bson:"group,omitempty"`
}

// Add is used to add a question.
func (question *Question) Add(db database.Database) (bool, error) {
	inserted, err := db.Insert("question", question)
	return inserted, err
}

// Get is used to get a question.
func (question *Question) Get(db database.Database, filter map[string]string) error {
	err := db.Get("question", filter, question)
	return err
}

// Update is used to update a question.
func (question *Question) Update(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.Update("question", filter, question)
	return updated, err
}

// QuestionGroup allows us to group similar questions.
type QuestionGroup struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
}

// Questions is a list of questions.
type Questions []Question

// Add method adds (a list of) questions at once.
func (questions *Questions) Add(db database.Database) error {
	err := db.InsertMultiple("question", questions)
	return err
}

// Get method gets all the questions in the question collection.
func (questions *Questions) Get(db database.Database, filter map[string]string) error {
	err := db.GetAll("question", filter, questions)
	return err
}

// Delete is used to delete a list of questions that match a filter.
func (questions Questions) Delete(db database.Database, filter map[string]string) (bool, error) {
	deleted, err := db.DeleteAll("question", filter)
	return deleted, err
}

// ToInterface converts questions (list of questions) into a list of interfaces, required by GetAll MongoDB.
func (questions Questions) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(questions))
	for i, item := range questions {
		interfaceObject[i] = item
	}
	return interfaceObject
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

// UpdateQuestion is used to add/remove translations to an existing question.
type UpdateQuestion map[string]interface{}

// Add is used to add a translation to an existing question.
func (question *UpdateQuestion) Add(db database.Database, filter map[string]string) (bool, error) {
	inserted, err := db.UpdateObject("question", filter, question)
	return inserted, err
}

// Remove is used to remove a translation to an existing question.
func (question *UpdateQuestion) Remove(db database.Database, filter map[string]string) (bool, error) {
	inserted, err := db.RemoveObject("question", filter, question)
	return inserted, err
}
