package models

// Question is the data for questions related to game types.
type Question struct {
	Rounds *Rounds `bson:"rounds" json:"rounds"`
}
