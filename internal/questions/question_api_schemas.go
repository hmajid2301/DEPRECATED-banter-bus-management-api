package questions

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal"
)

type QuestionIn struct {
	Content      string              `json:"content"                 validate:"required"`
	LanguageCode string              `json:"language_code,omitempty"                     default:"en"`
	Round        string              `json:"round,omitempty"`
	Group        *QuestionGroupInOut `json:"group,omitempty"`
}

type QuestionOut struct {
	Content string `json:"content" validate:"required"`
	Type    string `json:"type"                        enum:"answer,question"`
}

type QuestionGenericOut struct {
	Content string              `json:"content"         validate:"required"`
	Round   string              `json:"round,omitempty"`
	Enabled bool                `json:"enabled"`
	Group   *QuestionGroupInOut `json:"group,omitempty"`
}

type AllQuestionOut struct {
	IDs    []string `json:"ids"`
	Cursor string   `json:"cursor"`
}

type QuestionGroupInOut struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type"                     enum:"question,answer"`
}

type QuestionTranslationIn struct {
	Content string `json:"content" validate:"required"`
}

type AddQuestionInput struct {
	internal.GameParams
	QuestionIn
}

type LanguageParams struct {
	Language string `path:"language"`
}

type LanguageQueryParams struct {
	Language string `default:"en" query:"language"`
}

type QuestionIDParams struct {
	ID string `path:"question_id"`
}

type GroupNameParams struct {
	GroupName string `query:"group_name"`
}

type LimitParams struct {
	Limit int64 `query:"limit" default:"5" validate:"gte=0,lte=100"`
}

type QuestionInput struct {
	internal.GameParams
	QuestionIDParams
}

type QuestionLanguageInput struct {
	internal.GameParams
	LanguageParams
	QuestionIDParams
}

type GetQuestionIDsInput struct {
	internal.GameParams
	LimitParams
	Cursor string `query:"cursor"`
}

type GetQuestionInput struct {
	internal.GameParams
	LanguageParams
	QuestionIDParams
}

type AddTranslationInput struct {
	internal.GameParams
	LanguageParams
	QuestionIDParams
	QuestionTranslationIn
}

type GroupInput struct {
	internal.GameParams
	internal.RoundParams
}

type ListQuestionParams struct {
	internal.GameParams
	internal.RoundParams
	LanguageQueryParams
	GroupNameParams
	LimitParams
	Enabled string `query:"enabled" default:"all" enum:"enabled,disabled,all"`
	Random  bool   `query:"random"`
}
