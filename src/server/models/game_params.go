package models

type GameParams struct {
	Name string `bson:"name" json:"name" path:"name"`
}
