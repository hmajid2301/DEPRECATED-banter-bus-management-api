package models

// Question is the data for questions related to game types.
type Question struct {
	Content map[string]string `bson:"content"`
	Enabled *bool             `bson:"enabled,omitempty"`
}

// QuestionPool struct needed for user struct, but not implemented until question refactor
type QuestionPool struct {
	Unimplemented *bool `json:"unimplemented"`
}
