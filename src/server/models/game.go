package models

// Game is a specific game type. Includes details such as questions that can be used by the game.
type Game struct {
	Name     string `json:"name"`
	RulesURL string `json:"rules_url"`
	Enabled  bool   `json:"enabled"`
}
