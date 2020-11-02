package models

// NewQuestion is the data for new questions being added to games.
type NewQuestion struct {
	Content string `bson:"content" json:"content" validate:"required"`
	Round   string `bson:"round" json:"round" validate:"required,oneof=one two three"`
}
