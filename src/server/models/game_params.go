package models

// GameParams is the name of an existing game type.
type GameParams struct {
	Name string `json:"name" description:"The name of the game type" example:"quibly" path:"name"`
}
