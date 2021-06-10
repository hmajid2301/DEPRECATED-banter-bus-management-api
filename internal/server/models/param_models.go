package serverModels

// ListGameParams is the params from request.
type ListGameParams struct {
	Games string `query:"games" enum:"enabled,disabled,all" default:"all"`
}

// ListUserParams is the params to use in a GetAllUsers request.
type ListUserParams struct {
	AdminStatus string `query:"admin_status" url:"admin_status" enum:"admin,non_admin,all"        default:"all"`
	Privacy     string `query:"privacy"      url:"privacy"      enum:"private,friends,public,all" default:"all"`
	Membership  string `query:"membership"   url:"membership"   enum:"free,paid,all"              default:"all"`
}

// GameParams is the name of an existing game.
type GameParams struct {
	Name string `json:"name" description:"The name of the game." example:"quibly" path:"name"`
}

// LanguageParams is the language code for adding/removing questions.
type LanguageParams struct {
	Language string `json:"language" description:"The language code for the new question." example:"fr" path:"language"`
}

// StoryIDParams is the ID of the story in the database.
type StoryIDParams struct {
	StoryID string `description:"The id for the story." example:"2b45f6c6-d8be-4d13-9fc6-2f821c925774" path:"story_id"`
}
