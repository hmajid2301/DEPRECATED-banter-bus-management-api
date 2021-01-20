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
// More efficient way of storing strings for lookup than a slice hence it's a map
func (gameInfo *GameInfo) HasGroups(round string) bool {
	var gameRoundsWithGroups = map[string]struct{}{
		"fibbing_it.opinion":   {},
		"fibbing_it.free_form": {},
	}

	queryString := fmt.Sprintf("%s.%s", gameInfo.Name, round)
	_, isPresent := gameRoundsWithGroups[queryString]
	return isPresent
}

// Game is the interface for all game types.
type Game interface {
	GetInfo(rulesURL string) GameInfo
	GetQuestionPath(question GenericQuestion) string
	ValidateQuestionInput(question GenericQuestion) error
	GetQuestionPool() interface{}
	QuestionPoolToGenericQuestions(questions interface{}) ([]GenericQuestion, error)
}
