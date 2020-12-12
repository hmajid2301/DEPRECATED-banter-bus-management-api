package serverModels

// GenericQuestionGroup provides extra context to a question, when it belong to a group.
type GenericQuestionGroup struct {
	Name string `json:"name" description:"The name of the question group." example:"horse_group"`
	Type string `json:"type" description:"The type of the content."        example:"questions"   enum:"answers,questions"`
}
