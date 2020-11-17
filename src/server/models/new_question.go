package models

// NewQuestion is the data for new questions being added to games.
type NewQuestion struct {
	Content string `json:"content" validate:"required"`
	Round   string `json:"round"   validate:"required,oneof=one two three"`
}
