package serverModels

// QuestionTranslation is the data required to add an existing question in another language.
type QuestionTranslation struct {
	OriginalQuestion ReceiveQuestion        `json:"original_question" validate:"required"`
	NewQuestion      NewQuestionTranslation `json:"new_question"      validate:"required"`
}
