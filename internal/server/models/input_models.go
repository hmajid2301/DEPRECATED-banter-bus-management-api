package serverModels

// AddQuestionInput is used to add a new question to a game.
type AddQuestionInput struct {
	GameParams
	NewQuestion
}

// QuestionInput is the body data and params combined into a single struct.
type QuestionInput struct {
	GameParams
	LanguageParams
	QuestionIDParams
}

// AddTranslationInput is the body data and params combined into a single struct. Data required to add a new question translation.
type AddTranslationInput struct {
	GameParams
	LanguageParams
	QuestionIDParams
	QuestionTranslation
}

// GroupInput is the data for getting all groups from a certain round of a certain game.
type GroupInput struct {
	GameParams
	Round string `json:"round" url:"round" validate:"required" query:"round"`
}
