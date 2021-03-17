package models

import (
	"fmt"

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
	Name     string `bson:"name"`
	RulesURL string `bson:"rules_url,omitempty" json:"rules_url,omitempty"`
	Enabled  *bool  `bson:"enabled,omitempty"`
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

// Games is a list of game(s).
type Games []Game

// Add is used to add a list of games of the database.
func (games *Games) Add(db database.Database) error {
	err := db.InsertMultiple("game", games)
	return err
}

// Get is used to get a list of games.
func (games *Games) Get(db database.Database, filter map[string]string) error {
	err := db.GetAll("game", filter, games)
	return err
}

// Delete is used to delete a list of games that match a filter.
func (games Games) Delete(db database.Database, filter map[string]string) (bool, error) {
	deleted, err := db.DeleteAll("game", filter)
	return deleted, err
}

// ToInterface converts users (list of users) into a list of interfaces, required by GetAll MongoDB.
func (games Games) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(games))
	for i, item := range games {
		interfaceObject[i] = item
	}
	return interfaceObject
}
