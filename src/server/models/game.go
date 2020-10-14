package models

type Game struct {
	Name      string     `bson:"name" json:"name"`
	Questions []Question `bson:"questions" json:"questions"`
	RulesURL  string     `bson:"rules_url" json:"rules_url"`
	Enabled   bool       `bson:"enabled" json:"enabled"`
}
