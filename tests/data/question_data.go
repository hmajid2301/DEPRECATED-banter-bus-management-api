package data

import (
	"net/http"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
)

// AddQuestion is the test data for add questions to a game
var AddQuestion = []struct {
	TestDescription string
	Game            string
	Payload         interface{}
	Expected        int
}{
	{
		"Add a question to quibly and to round pair",
		"quibly",
		&serverModels.NewQuestion{
			Content: "this is another question?",
			Round:   "pair",
		}, http.StatusCreated,
	},
	{
		"Add a question to quibly and to round answer, language de",
		"quibly",
		&serverModels.NewQuestion{
			Content:      "what is the funniest thing ever told?",
			LanguageCode: "de",
			Round:        "answers",
		}, http.StatusCreated,
	},
	{
		"Add a question to quibly and to round group",
		"quibly",
		&serverModels.NewQuestion{
			Content: "what does ATGM stand for?",
			Round:   "group",
		}, http.StatusCreated,
	},
	{
		"Add a question to drawlosseum, language ur",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content:      "camel",
			LanguageCode: "ur",
		}, http.StatusCreated,
	},
	{
		"Add another question to drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "pencil",
		}, http.StatusCreated,
	},
	{
		"Add yet another question to drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "food fight",
		}, http.StatusCreated,
	},
	{
		"Add a question to fibbing it, round opinion new group bike group, language en",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content:      "do you love bikes?",
			LanguageCode: "en",
			Round:        "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "questions",
			},
		}, http.StatusCreated,
	},
	{
		"Add another question to fibbing it, round opinion new group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "how much does liam love bikes?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "questions",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round opinion existing group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "super love",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "answers",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round free_form existing group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What is love?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round free_form new group horse group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What is the fastest horse?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "horse_group",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round free_form existing group horse group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What is the second horse called?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "horse_group",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to never eat a vegetable again?",
			Round:   "likely",
		}, http.StatusCreated,
	},
	{
		"Add question to quibly, invalid round",
		"quibly",
		&serverModels.NewQuestion{
			Content: "This is another question?",
			Round:   "invalid",
		}, http.StatusBadRequest,
	},
	{
		"Add question to quibly, invalid2 round",
		"quibly",
		&serverModels.NewQuestion{
			Content: "This is another question?",
			Round:   "invalid2",
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, invalid round free_form",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What is the fastest horse?",
			Round:   "invalid_free_form",
			Group: &serverModels.Group{
				Name: "horse_group",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, invalid language code",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content:      "What is the fastest horse?",
			LanguageCode: "deed",
			Round:        "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "answers",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, round opinion invalid type",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "super love",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "answer",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, round opinion invalid answer type",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "super love",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "answer",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, round opinion invalid question type",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "super love",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "question",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, round opinion invalid type",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "super love",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "type",
			},
		}, http.StatusBadRequest,
	},
	{
		"game does not exist but missing content",
		"quibly v3",
		&serverModels.NewQuestion{}, http.StatusBadRequest,
	},
	{
		"game does not exist",
		"quibly_v2",
		&serverModels.NewQuestion{
			Content: "super love",
		}, http.StatusNotFound,
	},
	{
		"another game does not exist",
		"quibly v3",
		&serverModels.NewQuestion{
			Content: "super love",
		}, http.StatusNotFound,
	},
	{
		"Add a question that already exists to quibly and to round pair",
		"quibly",
		&serverModels.NewQuestion{
			Content: "this is another question?",
			Round:   "pair",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to quibly and to round answer",
		"quibly",
		&serverModels.NewQuestion{
			Content: "pink mustard",
			Round:   "answers",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to quibly and to round answer",
		"quibly",
		&serverModels.NewQuestion{
			Content:      "german",
			LanguageCode: "de",
			Round:        "answers",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to quibly and to round group",
		"quibly",
		&serverModels.NewQuestion{
			Content: "what does ATGM stand for?",
			Round:   "group",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "horse",
		}, http.StatusConflict,
	},
	{
		"Add another question that already exists to drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "pencil",
		}, http.StatusConflict,
	},
	{
		"Add yet another question that already exists to drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "food fight",
		}, http.StatusConflict,
	},
	{
		"Add a question to fibbing it that already exists, round opinion new group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "do you love bikes?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "questions",
			},
		}, http.StatusConflict,
	},
	{
		"Add another question to fibbing it that already exists, round opinion new group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "how much does liam love bikes?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "questions",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round opinion existing group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "super love",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "bike_group",
				Type: "answers",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round free_form existing group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What is love?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round free_form new group horse group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What is the fastest horse?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "horse_group",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round free_form existing group horse group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What is the second horse called?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "horse_group",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it tthat already exists, round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to never eat a vegetable again?",
			Round:   "likely",
		}, http.StatusConflict,
	},
	{
		"Add a question to quibly and to round pair",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What do you think about horses?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "questions",
			},
		}, http.StatusConflict,
	},
}

// UpdateQuestion is the test data for adding translations to questions.
var UpdateQuestion = []struct {
	TestDescription string
	Game            string
	Payload         interface{}
	Expected        int
}{
	{
		"Update question in quibly and round pair, new language fr",
		"quibly",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content:      "this is a question?",
				LanguageCode: "de",
				Round:        "pair",
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "this is a question?",
				LanguageCode: "fr",
			},
		},
		http.StatusOK,
	},
	{
		"Update question in quibly and round pair, replace exitsing language de",
		"quibly",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "pink mustard",
				Round:   "answers",
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "le german?",
				LanguageCode: "de",
			},
		},
		http.StatusOK,
	},
	{
		"Update question in quibly and round group, add new language de",
		"quibly",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content:      "this is a another question?",
				LanguageCode: "fr",
				Round:        "group",
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "Das ist eine andere Frage?",
				LanguageCode: "de",
			},
		},
		http.StatusOK,
	},
	{
		"Update question in quibly and round group, add another new language ur",
		"quibly",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content:      "this is a another question?",
				LanguageCode: "fr",
				Round:        "group",
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "Urdu question? Who knows?",
				LanguageCode: "ur",
			},
		},
		http.StatusOK,
	},
	{
		"Update question in drawlosseum",
		"drawlosseum",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "horse",
				Round:   "group",
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "ऊंट",
				LanguageCode: "hi",
			},
		},
		http.StatusOK,
	},
	{
		"Update question in drawlosseum, specify en (this should be default)",
		"drawlosseum",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content:      "spoon",
				LanguageCode: "en",
				Round:        "group",
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "spoon",
				LanguageCode: "hi",
			},
		},
		http.StatusOK,
	},
	{
		"Update question in fibbing it, round opinion",
		"fibbing_it",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "What do you think about horses?",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "horse_group",
					Type: "questions",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "Cosa ne pensi dei cavalli?",
				LanguageCode: "it",
			},
		},
		http.StatusOK,
	},
	{
		"Update question in fibbing it, round opinion and answers section",
		"fibbing_it",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "cool",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "horse_group",
					Type: "answers",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "Liebe",
				LanguageCode: "de",
			},
		}, http.StatusOK,
	},
	{
		"Update question in fibbing it, round free_form, language fr",
		"fibbing_it",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "Favourite bike colour?",
				Round:   "free_form",
				Group: &serverModels.Group{
					Name: "bike_group",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "was ist Liebe?",
				LanguageCode: "de",
			},
		}, http.StatusOK,
	},
	{
		"Update question in quibly, invalid round",
		"quibly",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "A question?",
				Round:   "invalid",
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "was ist Liebe?",
				LanguageCode: "de",
			},
		}, http.StatusBadRequest,
	},
	{
		"Update question in fibbing it, invalid round",
		"fibbing_it",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "Favourite bike colour?",
				Round:   "free_form2",
				Group: &serverModels.Group{
					Name: "bike_group",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "was ist Liebe?",
				LanguageCode: "de",
			},
		}, http.StatusBadRequest,
	},
	{
		"Update question in fibbing it, invalid group type answer (should be answers)",
		"fibbing_it",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "Favourite bike colour?",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "bike_group",
					Type: "answer",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "was ist Liebe?",
				LanguageCode: "de",
			},
		}, http.StatusBadRequest,
	},
	{
		"Missing content",
		"quibly",
		&serverModels.NewQuestion{}, http.StatusBadRequest,
	},
	{
		"Update question in fibbing it but invalid language code",
		"fibbing_it",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "Favourite bike colour?",
				Round:   "opinion",
				Group: &serverModels.Group{
					Name: "bike_group",
					Type: "answers",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "was ist Liebe?",
				LanguageCode: "ittt",
			},
		}, http.StatusBadRequest,
	},
	{
		"game does not exist",
		"quibly v3",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "Favourite bike colour?",
				Round:   "free_form",
				Group: &serverModels.Group{
					Name: "bike_group",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "was ist Liebe?",
				LanguageCode: "de",
			},
		}, http.StatusNotFound,
	},
	{
		"Original question doesn't exist",
		"fibbing_it",
		&serverModels.QuestionTranslation{
			OriginalQuestion: serverModels.NewQuestion{
				Content: "Favourite horse colour?",
				Round:   "free_form",
				Group: &serverModels.Group{
					Name: "bike_group",
				},
			},
			NewQuestion: serverModels.NewQuestionTranslation{
				Content:      "was ist Liebe?",
				LanguageCode: "de",
			},
		}, http.StatusNotFound,
	},
}

// RemoveQuestion is the test data for removing questions from game.
var RemoveQuestion = []struct {
	TestDescription string
	Game            string
	Payload         interface{}
	Expected        int
}{
	{
		"Delete a question quibly from round pair",
		"quibly",
		&serverModels.NewQuestion{
			Content: "this is a question?",
			Round:   "pair",
		}, http.StatusOK,
	},
	{
		"Delete a question quibly from round pair, language ur",
		"quibly",
		&serverModels.NewQuestion{
			Content:      "this is a question?",
			LanguageCode: "ur",
			Round:        "pair",
		}, http.StatusOK,
	},
	{
		"Delete a question quibly from round answers",
		"quibly",
		&serverModels.NewQuestion{
			Content: "pink mustard",
			Round:   "answers",
		}, http.StatusOK,
	},
	{
		"Delete a question quibly from round group, language fr",
		"quibly",
		&serverModels.NewQuestion{
			Content:      "this is a another question?",
			LanguageCode: "fr",
			Round:        "group",
		}, http.StatusOK,
	},
	{
		"Delete a question drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "horse",
		}, http.StatusOK,
	},
	{
		"Delete another question drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "spoon",
		}, http.StatusOK,
	},
	{
		"Delete a question to fibbing it, round opinion from group horse group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What do you think about horses?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "questions",
			},
		}, http.StatusOK,
	},
	{
		"Delete a answer to fibbing it, round opinion from group horse group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "cool",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "answers",
			},
		}, http.StatusOK,
	},
	{
		"Delete a answer to fibbing it, round free_form from group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "Favourite bike colour?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		}, http.StatusOK,
	},
	{
		"Delete a answer to fibbing it, round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to get arrested",
			Round:   "likely",
		}, http.StatusOK,
	},
	{
		"Delete another answer to fibbing it, round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
		}, http.StatusOK,
	},
	{
		"Delete a question quibly from round invalid",
		"quibly",
		&serverModels.NewQuestion{
			Content: "this is a question?",
			Round:   "invalid",
		}, http.StatusBadRequest,
	},
	{
		"Delete a question quibly from round content missing",
		"quibly",
		&serverModels.NewQuestion{
			Round: "group",
		}, http.StatusBadRequest,
	},
	{
		"Delete a question quibly from round pair that was already deleted",
		"quibly",
		&serverModels.NewQuestion{
			Content: "this is a question?",
			Round:   "pair",
		}, http.StatusNotFound,
	},
	{
		"Delete a question drawlosseum that was already deleted",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "horse",
		}, http.StatusNotFound,
	},
	{
		"Delete a question already removed from fibbing it, round free_form from group bike group",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "Favourite bike colour?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		}, http.StatusNotFound,
	},
	{
		"Delete a question already removed from fibbing it, round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to get arrested",
			Round:   "likely",
		}, http.StatusNotFound,
	},
	{
		"Delete another  already removed from fibbing it, round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
		}, http.StatusNotFound,
	},
}

// EnableQuestion test data used to test enable endpoint
var EnableQuestion = []struct {
	TestDescription string
	Game            string
	Payload         interface{}
	Expected        int
}{
	{
		"Enable a question, quibly and round pair",
		"quibly",
		&serverModels.NewQuestion{
			Content: "this is a question?",
			Round:   "pair",
		}, http.StatusOK,
	},
	{
		"Enable a question, quibly and round answers",
		"quibly",
		&serverModels.NewQuestion{
			Content:      "this is a another question?",
			LanguageCode: "fr",
			Round:        "group",
		}, http.StatusOK,
	},
	{
		"Enable a question, fibbing_it and round opinion",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What do you think about camels?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "questions",
			},
		}, http.StatusOK,
	},
	{
		"Enable an answer, fibbing_it and round opinion",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "cool",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "answers",
			},
		}, http.StatusOK,
	},
	{
		"Enable a question, fibbing_it and round free_form",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "Favourite bike colour?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		}, http.StatusOK,
	},
	{
		"Enable a question, fibbing_it and round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to get arrested",
			Round:   "likely",
		}, http.StatusOK,
	},
	{
		"Enable a question, drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "spoon",
		}, http.StatusOK,
	},
	{
		"Enable an already enabled question, drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "spoon",
		}, http.StatusOK,
	},
	{
		"Bad request invalid round, fibbing_it",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "spoon",
			Round:   "likely2",
		}, http.StatusBadRequest,
	},
	{
		"Bad request invalid content, fibbing_it",
		"fibbing_it",
		&serverModels.NewQuestion{}, http.StatusBadRequest,
	},
	{
		"Game does not exist",
		"quibly v3",
		&serverModels.NewQuestion{
			Content: "super love",
		}, http.StatusNotFound,
	},
}

// DisableQuestion test data used to test disable endpoint
var DisableQuestion = []struct {
	TestDescription string
	Game            string
	Payload         interface{}
	Expected        int
}{
	{
		"Disable a question, quibly and round pair",
		"quibly",
		&serverModels.NewQuestion{
			Content: "this is a question?",
			Round:   "pair",
		}, http.StatusOK,
	},
	{
		"Disable a question, quibly and round answers",
		"quibly",
		&serverModels.NewQuestion{
			Content: "pink mustard",
			Round:   "answers",
		}, http.StatusOK,
	},
	{
		"Disable a question, fibbing_it and round opinion",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "What do you think about camels?",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "questions",
			},
		}, http.StatusOK,
	},
	{
		"Disable an answer, fibbing_it and round opinion",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "lame",
			Round:   "opinion",
			Group: &serverModels.Group{
				Name: "horse_group",
				Type: "answers",
			},
		}, http.StatusOK,
	},
	{
		"Disable anquestion, fibbing_it and round free_form",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "Favourite bike colour?",
			Round:   "free_form",
			Group: &serverModels.Group{
				Name: "bike_group",
			},
		}, http.StatusOK,
	},
	{
		"Disable a question, fibbing_it and round likely",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "to eat ice-cream from the tub",
			Round:   "likely",
		}, http.StatusOK,
	},
	{
		"Disable a question, drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "spoon",
		}, http.StatusOK,
	},
	{
		"Disable a question, thats disabled drawlosseum",
		"drawlosseum",
		&serverModels.NewQuestion{
			Content: "spoon",
		}, http.StatusOK,
	},
	{
		"Bad request invalid round, fibbing_it",
		"fibbing_it",
		&serverModels.NewQuestion{
			Content: "spoon",
			Round:   "likely2",
		}, http.StatusBadRequest,
	},
	{
		"Bad request invalid content, fibbing_it",
		"fibbing_it",
		&serverModels.NewQuestion{}, http.StatusBadRequest,
	},
	{
		"Game does not exist",
		"quibly v3",
		&serverModels.NewQuestion{
			Content: "super love",
		}, http.StatusNotFound,
	},
}

// GetAllGroups is the data for the get groups tests
var GetAllGroups = []struct {
	TestDescription string
	Payload         *serverModels.GroupInput
	ExpectedGroups  []string
	ExpectedCode    int
}{
	{
		"Get all groups from questions from the opinion round in the Fibbing It game",
		&serverModels.GroupInput{
			GameParams: serverModels.GameParams{
				Name: "fibbing_it",
			},
			Round: "opinion",
		},
		[]string{
			"horse_group",
		},
		http.StatusOK,
	},

	{
		"Get all groups from questions from the free form round in the Fibbing It game",
		&serverModels.GroupInput{
			GameParams: serverModels.GameParams{
				Name: "fibbing_it",
			},
			Round: "free_form",
		},
		[]string{
			"bike_group",
			"cat_group",
		},
		http.StatusOK,
	},

	{
		"Try to get groups from a round in Fibbing It that does not have groups",
		&serverModels.GroupInput{
			GameParams: serverModels.GameParams{
				Name: "fibbing_it",
			},
			Round: "likely",
		},
		[]string{},
		http.StatusNotFound,
	},

	{
		"Try to get groups from a non-existent round",
		&serverModels.GroupInput{
			GameParams: serverModels.GameParams{
				Name: "fibbing_it",
			},
			Round: "genocide",
		},
		[]string{},
		http.StatusNotFound,
	},

	{
		"Try to get groups from a game that does not have groups",
		&serverModels.GroupInput{
			GameParams: serverModels.GameParams{
				Name: "quibly",
			},
			Round: "opinion",
		},
		[]string{},
		http.StatusNotFound,
	},
}
