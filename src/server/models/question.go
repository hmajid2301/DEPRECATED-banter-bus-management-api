package models

// Question is the data for questions related to game types.
type Question struct {
	One   []string `bson:"one" json:"one,omitempty"`
	Two   []string `bson:"two" json:"two,omitempty"`
	Three []string `bson:"three" json:"three,omitempty"`
}
