package models

// Rounds is the questions related to game types.
type Rounds struct {
	One   []string `json:"one"   validate:"omitempty"`
	Two   []string `json:"two"   validate:"omitempty"`
	Three []string `json:"three" validate:"omitempty"`
}
