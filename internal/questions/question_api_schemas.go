package questions

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal"
)

type QuestionIn struct {
	Content      string              `json:"content"                 description:"The question to add to a specific game."                  example:"This is a funny question?" validate:"required"`
	LanguageCode string              `json:"language_code,omitempty" description:"The language code for the question."                      example:"en"                                            default:"en"`
	Round        string              `json:"round,omitempty"         description:"If the game has rounds, specify the round in this field." example:"opinion"`
	Group        *QuestionGroupInOut `json:"group,omitempty"`
}

type QuestionOut struct {
	Content string `json:"content" description:"The question to add to a specific game." example:"This is a funny question?" validate:"required"`
	Type    string `json:"type"    description:"The type of content question or answer."                                                         enum:"answer,question"`
}

type QuestionGenericOut struct {
	Content string              `json:"content"         description:"The question to add to a specific game."                                example:"This is a funny question?" validate:"required"`
	Round   string              `json:"round,omitempty" description:"If the game has rounds, specify the round in this field."               example:"opinion"`
	Enabled bool                `json:"enabled"         description:"True if the question is enabled and can be used in a game, else false."`
	Group   *QuestionGroupInOut `json:"group,omitempty"`
}

type AllQuestionOut struct {
	IDs    []string `json:"ids"    description:"All the question ids."`
	Cursor string   `json:"cursor" description:"The next question id (for pagination)."`
}

type QuestionGroupInOut struct {
	Name string `json:"name" description:"The name of the group."         example:"animal_group" validate:"required"`
	Type string `json:"type" description:"The type of the content group." example:"question"                         enum:"question,answer"`
}

type QuestionTranslationIn struct {
	Content string `json:"content" description:"The question in the new language" example:"Willst du eine Frage?" validate:"required"`
}

type AddQuestionInput struct {
	internal.GameParams
	QuestionIn
}

type LanguageParams struct {
	Language string `description:"The language code for the new question." example:"fr" path:"language"`
}

type LanguageQueryParams struct {
	Language string `description:"The language code for the new question." example:"fr" default:"en" query:"language"`
}

type QuestionIDParams struct {
	ID string `description:"The id for a specific question." example:"a-random-id" path:"question_id"`
}

type GroupNameParams struct {
	GroupName string `description:"The name of the group." example:"horse" query:"group_name"`
}

type LimitParams struct {
	Limit int64 `description:"The number of questions to retrieve." query:"limit" default:"5" validate:"gte=0,lte=100"`
}

type QuestionInput struct {
	internal.GameParams
	LanguageParams
	QuestionIDParams
}

type GetQuestionIDsInput struct {
	internal.GameParams
	LimitParams
	Cursor string `description:"The ID to start at for retrieving ID" example:"60e777f2d24d7d711e971aee" query:"cursor"`
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
	Enabled string `description:"If set to false will retrieve questions that are not enabled." query:"enabled" default:"all" enum:"enabled,disabled,all"`
	Random  bool   `description:"If set will retrieve questions randomly."                      query:"random"`
}
