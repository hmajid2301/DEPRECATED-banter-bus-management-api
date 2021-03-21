package serverModels

// QuestionInput is the body data and params combined into a single struct.
type QuestionInput struct {
	GameParams
	LanguageParams
	NewQuestion
}

// AddTranslationInput is the body data and params combined into a single struct. Data required to add a new question translation.
type AddTranslationInput struct {
	GameParams
	LanguageParams
	QuestionTranslation
}

// GroupInput is the data for getting all groups from a certain round of a certain game.
type GroupInput struct {
	GameParams
	Round string `json:"round" url:"round" validate:"required" query:"round"`
}

// QuestionPoolInput is the combined data to create a new question pool.
type QuestionPoolInput struct {
	UserParams
	Pool
}

// UpdateQuestionPoolInput is the combined data (params + body) required to update an existing question pool.
// Such as adding or removing a new question.
type UpdateQuestionPoolInput struct {
	UserParams
	PoolParams
	NewQuestion
}
