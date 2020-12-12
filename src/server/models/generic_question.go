package serverModels

// GenericQuestion is generic structure all questions can take, has all the required fields for any question.
type GenericQuestion struct {
	Content string                `json:"content"         description:"The question to add to a specific game."                  example:"This is a funny question?" validate:"required"`
	Round   string                `json:"round,omitempty" description:"If the game has rounds, specify the round in this field." example:"opinion"`
	Group   *GenericQuestionGroup `json:"group,omitempty"`
}
