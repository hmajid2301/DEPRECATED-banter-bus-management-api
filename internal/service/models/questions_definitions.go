package models

// StoryAnswerType is the interface for any story answer types added to a user.
type StoryAnswerType interface {
	NewAnswer()
}
