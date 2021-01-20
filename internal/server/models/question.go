package serverModels

// NewQuestion is the data for new questions being added to games.
type NewQuestion struct {
	Content      string `json:"content"                 description:"The question to add to a specific game."                  example:"This is a funny question?" validate:"required"`
	LanguageCode string `json:"language_code,omitempty" description:"The language code for the question."                      example:"en"                                            default:"en"`
	Round        string `json:"round,omitempty"         description:"If the game has rounds, specify the round in this field." example:"opinion"`
	Group        *Group `json:"group,omitempty"`
}

// GenericQuestion is generic structure all questions can take, has all the required fields for any question.
type GenericQuestion struct {
	Content string                `json:"content"         description:"The question to add to a specific game."                  example:"This is a funny question?" validate:"required"`
	Round   string                `json:"round,omitempty" description:"If the game has rounds, specify the round in this field." example:"opinion"`
	Group   *GenericQuestionGroup `json:"group,omitempty"`
}

// GenericQuestionGroup provides extra context to a question, when it belong to a group.
type GenericQuestionGroup struct {
	Name string `json:"name" description:"The name of the question group." example:"horse_group"`
	Type string `json:"type" description:"The type of the content."        example:"questions"   enum:"answers,questions"`
}

// Group is the data for new questions being added to some game types.
type Group struct {
	Name string `json:"name" description:"The name of the group."         example:"animal_group" validate:"required"`
	Type string `json:"type" description:"The type of the content group." example:"questions"                        enum:"questions,answers"`
}

// GroupInput is the data for getting all groups from a certain round of a certain game type
type GroupInput struct {
	GameName string `json:"game_name" url:"game_name" description:"The name of the game" example:"fibbing_it" validate:"required" path:"name"`
	Round    string `json:"round"     url:"round"                                                             validate:"required"             query:"round"`
}

// QuestionTranslation is the data required to add an existing question in another language.
type QuestionTranslation struct {
	OriginalQuestion NewQuestion            `json:"original_question" validate:"required"`
	NewQuestion      NewQuestionTranslation `json:"new_question"      validate:"required"`
}

// NewQuestionTranslation is the data required by the new question being added.
type NewQuestionTranslation struct {
	Content      string `json:"content"       description:"The question in the new language"        example:"Willst du eine Frage?" validate:"required"`
	LanguageCode string `json:"language_code" description:"The language code for the new question." example:"fr"                    validate:"required"`
}

// QuestionInput is the body data and params combined into a single struct.
type QuestionInput struct {
	GameParams
	NewQuestion
}

// UpdateQuestionInput is the body data and params combined into a single struct. Data required to update a question.
type UpdateQuestionInput struct {
	GameParams
	QuestionTranslation
}
