package models

// Question is the data for questions related to game types.
type Question struct {
	Content map[string]string `bson:"content"`
	Enabled *bool             `bson:"enabled,omitempty"`
}

// QuiblyQuestions is the data for questions related to the Quibly game.
type QuiblyQuestions struct {
	Pair    []Question `bson:"pair,omitempty"`
	Answers []Question `bson:"answers,omitempty"`
	Group   []Question `bson:"group,omitempty"`
}

// DrawlosseumQuestions is the data required to play the Drawlosseum game.
type DrawlosseumQuestions struct {
	Drawings []Question `bson:"drawings"`
}

// FibbingItQuestions is the data for questions related to the Fibbing It game.
// NOTE: Some fields are a map of questions, because questions will be grouped with other similar questions
type FibbingItQuestions struct {
	Opinion  map[string]map[string][]Question `bson:"opinion"`
	FreeForm map[string][]Question            `bson:"free_form"`
	Likely   []Question                       `bson:"likely"`
}

// GenericQuestion is generic structure all questions can take, has all the required fields for any question.
type GenericQuestion struct {
	Content      string
	Round        string
	LanguageCode string
	Group        *GenericQuestionGroup
}

// GenericQuestionGroup provides extra context to a question, when it belong to a group.
type GenericQuestionGroup struct {
	Name string
	Type string
}
