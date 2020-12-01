package serverModels

// GameParams is the name of an existing game.
type GameParams struct {
	Name string `json:"name" description:"The name of the game." example:"quibly" path:"name"`
}
