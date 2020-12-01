package serverModels

// NewUser struct is the data needed to add a new user via the API
type NewUser struct {
	Username   string `json:"username"   description:"The screen name of the user." example:"lmoz25" validate:"required"`
	Membership string `json:"membership" description:"membership for a user"        example:"free"   validate:"required,oneof=free paid"`
	Admin      *bool  `json:"admin"      description:"Whether the user is admin"                     validate:"required"`
}
