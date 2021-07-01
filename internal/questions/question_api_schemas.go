package questions

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal"
)

type QuestionIn struct {
	Content      string           `json:"content"                 description:"The question to add to a specific game."                  example:"This is a funny question?" validate:"required"`
	LanguageCode string           `json:"language_code,omitempty" description:"The language code for the question."                      example:"en"                                            default:"en"`
	Round        string           `json:"round,omitempty"         description:"If the game has rounds, specify the round in this field." example:"opinion"`
	Group        *QuestionGroupIn `json:"group,omitempty"`
}

type QuestionGroupIn struct {
	Name string `json:"name" description:"The name of the group."         example:"animal_group" validate:"required"`
	Type string `json:"type" description:"The type of the content group." example:"questions"                        enum:"questions,answers"`
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

type QuestionIDParams struct {
	ID string `description:"The id for a specific question." example:"a-random-id" path:"question_id"`
}

type QuestionInput struct {
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
	Round string `url:"round" validate:"required" query:"round"`
}
