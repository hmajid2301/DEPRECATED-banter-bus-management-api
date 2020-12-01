package models

// FibbingItQuestions is the data for questions related to the Fibbing It game.
// NOTE: Some fields are a map of questions, because questions will be grouped with other similar questions
type FibbingItQuestions struct {
	Opinion  []map[string]Question `bson:"opinion"`
	FreeForm []map[string]Question `bson:"free_form"`
	Likely   []Question            `bson:"likely"`
}
