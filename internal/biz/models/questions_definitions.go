package models

// QuestionPoolType is the interface for any question pool added to a user.
type QuestionPoolType interface {
	EmptyPoolQuestions()
}

// QuestionType is the interface for any question types added to a game.
type QuestionType interface {
	EmptyQuestions()
}

// StoryAnswerType is the interface for any story answer types added to a user.
type StoryAnswerType interface {
	EmptyAnswer()
}
