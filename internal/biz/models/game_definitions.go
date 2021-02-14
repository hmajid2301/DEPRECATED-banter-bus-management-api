package models

// IGame is the interface for game(s).
type IGame interface {
	NewGame(rulesURL string) Game
	GetQuestionPath(question GenericQuestion) string
	ValidateQuestionInput(question GenericQuestion) error
	NewQuestionPool() QuestionPoolType
	QuestionPoolToGenericQuestions(questions QuestionPoolType) ([]GenericQuestion, error)
}
