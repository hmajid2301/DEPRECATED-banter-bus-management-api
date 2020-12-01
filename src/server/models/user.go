package serverModels

// User struct to hold a database entry for a user, i.e a player with an account
type User struct {
	Username      string           `json:"username"       description:"The screen name of the user."               example:"lmoz25"`
	Admin         *bool            `json:"admin"          description:"Whether or not the user is admin."          example:"true"`
	Privacy       string           `json:"privacy"        description:"The privacy level of the user."             example:"public"`
	Membership    string           `json:"membership"     description:"The level of membership the user has."      example:"free"`
	Preferences   *UserPreferences `json:"preferences"    description:"Collection of user preferences."`
	Friends       []Friend         `json:"friends"        description:"List of friends the user has."`
	Stories       []Story          `json:"stories"        description:"Details of stories the user has."`
	QuestionPools []QuestionPool   `json:"question_pools" description:"List of question pools the user has saved."`
}

// UserPreferences struct to hold information on a user's preferences
type UserPreferences struct {
	LanguageCode string `json:"language_code" description:"Details of the user's preferred display and question language stored using ISO 2-letter language abbreviations" example:"en"`
}

// Friend struct to hold details about a user's friend
type Friend struct {
	Username string `json:"username" description:"The screen name of the friend" example:"seeb123"`
}

// UserParams is the username of an existing user
type UserParams struct {
	Username string `json:"username" description:"The screen name of the user" example:"lmoz25" path:"name"`
}

// QuestionPool struct needed for user struct, but not implemented until question refactor
type QuestionPool struct {
	Unimplemented *bool `json:"unimplemented"`
}
