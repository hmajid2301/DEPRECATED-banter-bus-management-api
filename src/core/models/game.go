package models

// Game is the data required by all game types. The only thing that varies between different games
// is how they store questions.
type Game struct {
	Name      string      `bson:"name"`
	RulesURL  string      `bson:"rules_url,omitempty"`
	Enabled   *bool       `bson:"enabled,omitempty"`
	Questions interface{} `bson:"questions,omitempty"`
}
