package models

// QuiblyQuestions is the data for questions related to the Quibly game.
type QuiblyQuestions struct {
	Pair   []Question `bson:"pair,omitempty"`
	Answer []Question `bson:"answer,omitempty"`
	Group  []Question `bson:"group,omitempty"`
}
