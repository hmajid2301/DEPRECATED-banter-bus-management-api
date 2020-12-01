package serverModels

// Group is the data for new questions being added to some game types.
type Group struct {
	Name string `json:"name" description:"The name of the group."         example:"animal_group" validate:"required"`
	Type string `json:"type" description:"The type of the content group." example:"questions"                        enum:"questions,answers"`
}
