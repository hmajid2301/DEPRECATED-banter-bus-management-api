package serverModels

// Group is the data for new questions being added to some game types.
type Group struct {
	Name string `json:"name" description:"The name of the group."         example:"animal_group" validate:"required"`
	Type string `json:"type" description:"The type of the content group." example:"questions"                        enum:"questions,answers"`
}

// GroupInput is the data for getting all groups from a certain round of a certain game type
type GroupInput struct {
	GameName string `json:"game_name" url:"game_name" description:"The name of the game" example:"fibbing_it" validate:"required" path:"name"`
	Round    string `json:"round" query:"round" url:"round" validate:"required"`
}
