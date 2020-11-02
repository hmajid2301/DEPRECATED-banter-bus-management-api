package models

// GameParams is the name of an existing game type.
type GameParams struct {
	Name string `bson:"name" json:"name" path:"name"`
}
