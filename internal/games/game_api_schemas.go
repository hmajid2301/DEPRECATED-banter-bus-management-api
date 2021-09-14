package games

import "gitlab.com/banter-bus/banter-bus-management-api/internal"

type GameIn struct {
	Name        string `json:"name"         validate:"required"`
	RulesURL    string `json:"rules_url"    validate:"required"`
	Description string `json:"description"  validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
}

type GameOut struct {
	Name        string `json:"name"`
	RulesURL    string `json:"rules_url"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
	DisplayName string `json:"display_name" validate:"required"`
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
