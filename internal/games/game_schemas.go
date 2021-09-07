package games

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

const QUIBLY = "quibly"
const FIBBINGIT = "fibbing_it"
const DRAWLOSSEUM = "drawlosseum"

type Game struct {
	Name        string `bson:"name"`
	RulesURL    string `bson:"rules_url,omitempty"   json:"rules_url,omitempty"`
	Enabled     *bool  `bson:"enabled,omitempty"`
	Description string `bson:"description,omitempty"`
}

func (game *Game) Add(db database.Database) (bool, error) {
	inserted, err := db.Insert("game", game)
	return inserted, err
}

func (game *Game) Get(db database.Database, filter map[string]interface{}) error {
	err := db.Get("game", filter, game)
	return err
}

func (game *Game) Update(db database.Database, filter map[string]interface{}) (bool, error) {
	updated, err := db.Update("game", filter, game)
	return updated, err
}

type Games []Game

func (games *Games) Add(db database.Database) error {
	err := db.InsertMultiple("game", games)
	return err
}

func (games *Games) Get(db database.Database, filter map[string]interface{}) error {
	err := db.GetAll("game", filter, games)
	return err
}

func (games *Games) GetWithLimit(db database.Database, filter map[string]interface{}, limit int64) error {
	return nil
}

func (games Games) Delete(db database.Database, filter map[string]interface{}) (bool, error) {
	deleted, err := db.DeleteAll("game", filter)
	return deleted, err
}

func (games Games) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(games))
	for i, item := range games {
		interfaceObject[i] = item
	}
	return interfaceObject
}
