package models

// Gamer is the interface for game(s).
type Gamer interface {
	GetQuestionPath(question GenericQuestion) string
	ValidateQuestion(question GenericQuestion) error
	NewQuestionPool() QuestionPoolType
	QuestionPoolToGenericQuestions(questions QuestionPoolType) ([]GenericQuestion, error)
}
