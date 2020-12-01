package models

// Question is the data for questions related to game types.
type Question struct {
	Content map[string]string `bson:"content"`
	Enabled *bool             `bson:"enabled,omitempty"`
}
