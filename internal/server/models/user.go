package serverModels

// NewUser struct is the data needed to add a new user via the API
type NewUser struct {
	Username   string `json:"username"   description:"The screen name of the user." example:"lmoz25" validate:"required"`
	Membership string `json:"membership" description:"membership for a user"        example:"free"   validate:"required,oneof=free paid"`
	Admin      *bool  `json:"admin"      description:"Whether the user is admin"                     validate:"required"`
}

// User struct to hold a database entry for a user, i.e a player with an account
type User struct {
	Username    string           `json:"username"    description:"The screen name of the user."          example:"lmoz25"`
	Admin       *bool            `json:"admin"       description:"Whether or not the user is admin."`
	Privacy     string           `json:"privacy"     description:"The privacy level of the user."        example:"public"`
	Membership  string           `json:"membership"  description:"The level of membership the user has." example:"free"`
	Preferences *UserPreferences `json:"preferences" description:"Collection of user preferences."`
	Friends     []Friend         `json:"friends"     description:"List of friends the user has."`
}

// UserPreferences struct to hold information on a user's preferences
type UserPreferences struct {
	LanguageCode string `json:"language_code" description:"Details of the user's preferred display and question language stored using ISO 2-letter language abbreviations" example:"en"`
}

// Friend struct to hold details about a user's friend
type Friend struct {
	Username string `json:"username" description:"The screen name of the friend" example:"seeb123"`
}

// QuestionPool is the list of custom questions for different game.
type QuestionPool struct {
	PoolName  string            `json:"pool_name" description:"The unique name of the question pool."      example:"my_pool"`
	GameName  string            `json:"game_name" description:"The type of game the story pertains to"     example:"quibly"`
	Privacy   string            `json:"privacy"   description:"The privacy setting for this question pool"                   enum:"public,private,friends"`
	Questions []GenericQuestion `json:"questions" description:"List of questions in this pool."`
}
