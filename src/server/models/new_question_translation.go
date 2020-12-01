package serverModels

// NewQuestionTranslation is the data required by the new question being added.
type NewQuestionTranslation struct {
	Content      string `json:"content"       description:"The question in the new language"        example:"Willst du eine Frage?" validate:"required"`
	LanguageCode string `json:"language_code" description:"The language code for the new question." example:"fr"                    validate:"required"`
}
