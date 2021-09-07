package games

import "gitlab.com/banter-bus/banter-bus-management-api/internal"

type GameIn struct {
	Name        string `json:"name"        description:"The name of the new game "         example:"quibly"              validate:"required"`
	RulesURL    string `json:"rules_url"   description:"The URL to the rules of the game." example:"gitlab.com/rules.md" validate:"required"`
	Description string `json:"description" description:"A brief description of the game."                                validate:"required"`
}

type GameOut struct {
	Name        string `json:"name"        description:"The name of the new game."           example:"quibly"`
	RulesURL    string `json:"rules_url"   description:"The URL to the rules of the game."   example:"gitlab.com/rules.md"`
	Enabled     bool   `json:"enabled"     description:"If set to true the game is enabled." example:"false"`
	Description string `json:"description" description:"A brief description of the game."`
}

type ListGameParams struct {
	Games string `query:"games" enum:"enabled,disabled,all" default:"all"`
}

type ListUserParams struct {
	AdminStatus string `query:"admin_status" enum:"admin,non_admin,all"        default:"all"`
	Privacy     string `query:"privacy"      enum:"private,friends,public,all" default:"all"`
	Membership  string `query:"membership"   enum:"free,paid,all"              default:"all"`
}

type GroupInput struct {
	internal.GameParams
	internal.RoundParams
}
