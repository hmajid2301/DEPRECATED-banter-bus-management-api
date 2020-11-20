package models

// NewGame is the data needed to add a new game type to the API.
type NewGame struct {
	Name     string `json:"name"      description:"The name of the new game type."                        example:"quibly"              validate:"required"`
	RulesURL string `json:"rules_url" description:"The URL to the rules of the game, as a markdown file." example:"gitlab.com/rules.md" validate:"required"`
}
