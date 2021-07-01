package questions

import "gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"

type Question struct {
	ID       string            `bson:"id"`
	GameName string            `bson:"game_name"       json:"game_name"`
	Round    string            `bson:"round,omitempty"`
	Enabled  *bool             `bson:"enabled"`
	Content  map[string]string `bson:"content"`
	Group    QuestionGroup     `bson:"group,omitempty"`
}

func (question *Question) Add(db database.Database) (bool, error) {
	inserted, err := db.Insert("question", question)
	return inserted, err
}

func (question *Question) Get(db database.Database, filter map[string]string) error {
	err := db.Get("question", filter, question)
	return err
}

func (question *Question) Update(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.Update("question", filter, question)
	return updated, err
}

type QuestionGroup struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
}

type Questions []Question

func (questions *Questions) Add(db database.Database) error {
	err := db.InsertMultiple("question", questions)
	return err
}

func (questions *Questions) Get(db database.Database, filter map[string]string) error {
	err := db.GetAll("question", filter, questions)
	return err
}

func (questions Questions) Delete(db database.Database, filter map[string]string) (bool, error) {
	deleted, err := db.DeleteAll("question", filter)
	return deleted, err
}

func (questions Questions) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(questions))
	for i, item := range questions {
		interfaceObject[i] = item
	}
	return interfaceObject
}

type GenericQuestion struct {
	Content      string
	Round        string
	LanguageCode string
	Group        *GenericQuestionGroup
}

type GenericQuestionGroup struct {
	Name string
	Type string
}

type UpdateQuestion map[string]interface{}

func (question *UpdateQuestion) Add(db database.Database, filter map[string]string) (bool, error) {
	inserted, err := db.UpdateObject("question", filter, question)
	return inserted, err
}

func (question *UpdateQuestion) Remove(db database.Database, filter map[string]string) (bool, error) {
	inserted, err := db.RemoveObject("question", filter, question)
	return inserted, err
}
