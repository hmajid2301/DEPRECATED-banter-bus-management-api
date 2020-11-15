package dbmodels

// Game is a specific game type. Includes details such as questions that can be used by the game.
type Game struct {
	Name      string    `bson:"name"`
	Questions *Question `bson:"questions,omitempty"`
	RulesURL  string    `bson:"rules_url,omitempty"`
	Enabled   *bool     `bson:"enabled,omitempty"`
}
