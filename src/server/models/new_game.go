package models

//NewGame is the data needed to add a new game type to the API.
type NewGame struct {
	Name     string `bson:"name" json:"name" validate:"required"`
	RulesURL string `bson:"rules_url" json:"rules_url" validate:"required"`
}
