package models

// Rounds is the questions related to game types.
type Rounds struct {
	One   []string `bson:"one" json:"one" validate:"omitempty"`
	Two   []string `bson:"two" json:"two" validate:"omitempty"`
	Three []string `bson:"three" json:"three" validate:"omitempty"`
}
