package serverModels

// NewGame is the data needed to add a new game to the API.
type NewGame struct {
	Name     string `json:"name"      description:"The name of the new game "         example:"quibly"              validate:"required"`
	RulesURL string `json:"rules_url" description:"The URL to the rules of the game." example:"gitlab.com/rules.md" validate:"required"`
}

// Game struct includes details such as questions that can be used by the game.
type Game struct {
	Name     string `json:"name"      description:"The name of the new game."           example:"quibly"`
	RulesURL string `json:"rules_url" description:"The URL to the rules of the game."   example:"gitlab.com/rules.md"`
	Type     string `json:"type"      description:"The type of the new game."           example:"quibly"              validate:"required,oneof=quibly fibbing_it drawlosseum"`
	Enabled  bool   `json:"enabled"   description:"If set to true the game is enabled." example:"false"`
}
