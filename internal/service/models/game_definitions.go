package models

// Gamer is the interface for game(s).
type Gamer interface {
	NewGame(rulesURL string) Game
	GetQuestionPath(question GenericQuestion) string
	ValidateQuestion(question GenericQuestion) error
	NewQuestionPool() QuestionPoolType
	QuestionPoolToGenericQuestions(questions QuestionPoolType) ([]GenericQuestion, error)
}
