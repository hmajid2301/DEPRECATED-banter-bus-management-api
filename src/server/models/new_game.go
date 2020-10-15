package models

type NewGame struct {
	Name     string `bson:"name" json:"name" validate:"required"`
	RulesURL string `bson:"rules_url" json:"rules_url" validate:"required"`
}
