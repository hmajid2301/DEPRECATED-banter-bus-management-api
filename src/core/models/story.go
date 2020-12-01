package models

// Story struct to contain information about a user story
type Story struct {
	GameName string        `bson:"gamename"`
	Question string        `bson:"question"`
	Answers  []StoryAnswer `bson:"answers"`
}

// StoryAnswer struct to contain information needed to store an answer for a player's story
type StoryAnswer struct {
	Answer string `bson:"answer"`
	Votes  int    `bson:"votes"`
}
