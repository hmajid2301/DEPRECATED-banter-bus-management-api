package models

// Game is a specific game type. Includes details such as questions that can be used by the game.
type Game struct {
	Name     string `json:"name"      description:"The name of the new game type."                        example:"quibly"`
	RulesURL string `json:"rules_url" description:"The URL to the rules of the game, as a markdown file." example:"gitlab.com/rules.md"`
	Enabled  bool   `json:"enabled"   description:"If set to true the game is enabled."                   example:"false"`
}
