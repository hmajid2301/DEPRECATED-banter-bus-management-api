package serverModels

// UpdateQuestionInput is the body data and params combined into a single struct. Data required to update a question.
type UpdateQuestionInput struct {
	GameParams
	QuestionTranslation
}
