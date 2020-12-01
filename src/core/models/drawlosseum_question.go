package models

// DrawlosseumQuestions is the data required to play the Drawlosseum game.
type DrawlosseumQuestions struct {
	Drawings []Question `bson:"drawings"`
}
