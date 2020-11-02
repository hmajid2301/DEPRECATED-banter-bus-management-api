package models

//Question is the data for questions related to game types.
type Question struct {
	Question string ` bson:"question" json:"question"`
}
