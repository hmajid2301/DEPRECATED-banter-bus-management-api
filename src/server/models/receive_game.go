package serverModels

// ReceiveGame is the data needed to add a new game to the API.
type ReceiveGame struct {
	Name     string `json:"name"      description:"The name of the new game "         example:"quibly"              validate:"required"`
	RulesURL string `json:"rules_url" description:"The URL to the rules of the game." example:"gitlab.com/rules.md" validate:"required"`
}
