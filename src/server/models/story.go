package serverModels

// Story struct to contain information about a user story
type Story struct {
	GameName string        `json:"game_name" description:"The type of game the story pertains to"                                example:"quibly"`
	Question string        `json:"question"  description:"The text of the question"                                              example:"how many fish?"`
	Answers  []StoryAnswer `json:"answers"   description:"The answer(s) given by the player(s) and the number of votes they got"`
}

// StoryAnswer struct to contain information needed to store an answer for a player's story
type StoryAnswer struct {
	Answer string `json:"answer" description:"The body of the answer, be it text or image" example:"4 fish"`
	Votes  int    `json:"votes"  description:"The number of votes the answer got"          example:"7"`
}
