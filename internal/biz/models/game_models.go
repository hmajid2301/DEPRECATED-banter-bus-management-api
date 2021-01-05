package models

import "fmt"

// GameInfo is the data required by all game types. The only thing that varies between different games
// is how they store questions.
type GameInfo struct {
	Name      string      `bson:"name"`
	RulesURL  string      `bson:"rules_url,omitempty" json:"rules_url,omitempty"`
	Enabled   *bool       `bson:"enabled,omitempty"`
	Questions interface{} `bson:"questions,omitempty"`
}

// HasGroups checks if the game has question groups for the specified round
func (gameInfo *GameInfo) HasGroups(round string) bool {
	queryString := fmt.Sprintf("%s.%s", gameInfo.Name, round)
	_, isPresent := gameRoundsWithGroups[queryString]
	return isPresent
}

// PlayableGame interface for all game types.
type PlayableGame interface {
	AddGame(string) (bool, error)
	GetQuestionPath() string
	ValidateQuestionInput() error
}

// More efficient way of storing strings for lookup than a slice
var gameRoundsWithGroups = map[string]struct{}{
	"fibbing_it.opinion":   {},
	"fibbing_it.free_form": {},
}
