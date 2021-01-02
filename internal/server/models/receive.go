package serverModels

// ReceiveQuestion is the data for new questions being added to games.
type ReceiveQuestion struct {
	Content      string `json:"content"                 description:"The question to add to a specific game."                  example:"This is a funny question?" validate:"required"`
	LanguageCode string `json:"language_code,omitempty" description:"The language code for the question."                      example:"en"                                            default:"en"`
	Round        string `json:"round,omitempty"         description:"If the game has rounds, specify the round in this field." example:"opinion"`
	Group        *Group `json:"group,omitempty"`
}

// ReceiveGame is the data needed to add a new game to the API.
type ReceiveGame struct {
	Name     string `json:"name"      description:"The name of the new game "         example:"quibly"              validate:"required"`
	RulesURL string `json:"rules_url" description:"The URL to the rules of the game." example:"gitlab.com/rules.md" validate:"required"`
}
