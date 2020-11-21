package models

// ReceiveQuestion is the data for new questions being added to games.
type ReceiveQuestion struct {
	Content string `json:"content" description:"The question to add to a specific game type." example:"This is a funny question?" validate:"required"`
	Round   string `json:"round"   description:"Which round to add the question to."          example:"one"                       validate:"required,oneof=one two three"`
}
