package models

// Game is a specific game type. Includes details such as questions that can be used by the game.
type Game struct {
	Name      string     `bson:"name" json:"name"`
	Questions []Question `bson:"questions" json:"questions"`
	RulesURL  string     `bson:"rules_url" json:"rules_url"`
	Enabled   bool       `bson:"enabled" json:"enabled"`
}
