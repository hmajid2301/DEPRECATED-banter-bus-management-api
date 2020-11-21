package models

// ListGameParams is the params from request.
type ListGameParams struct {
	Games string `query:"games" enum:"enabled,disabled,all" default:"all"`
}
