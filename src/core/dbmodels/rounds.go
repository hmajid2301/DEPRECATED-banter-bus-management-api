package dbmodels

// Rounds is the questions related to game types.
type Rounds struct {
	One   []string `bson:"one"`
	Two   []string `bson:"two"`
	Three []string `bson:"three"`
}
