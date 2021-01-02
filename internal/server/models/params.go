package serverModels

// ListGameParams is the params from request.
type ListGameParams struct {
	Games string `query:"games" enum:"enabled,disabled,all" default:"all"`
}

// ListUserParams is the params to use in a GetAllUsers request
type ListUserParams struct {
	AdminStatus string `query:"admin_status" url:"admin_status" enum:"admin,non_admin,all"        default:"all"`
	Privacy     string `query:"privacy"      url:"privacy"      enum:"private,friends,public,all" default:"all"`
	Membership  string `query:"membership"   url:"membership"   enum:"free,paid,all"              default:"all"`
}

// GameParams is the name of an existing game.
type GameParams struct {
	Name string `json:"name" description:"The name of the game." example:"quibly" path:"name"`
}

// UserParams is the username of an existing user
type UserParams struct {
	Username string `json:"username" description:"The screen name of the user" example:"lmoz25" path:"name"`
}
