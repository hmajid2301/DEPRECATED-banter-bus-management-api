package data

import (
	"net/http"

	"gitlab.com/banter-bus/banter-bus-management-api/internal"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/questions"
)

var AddQuestion = []struct {
	TestDescription string
	Game            string
	Payload         interface{}
	Expected        int
}{
	{
		"Add a question to quibly and to round pair",
		"quibly",
		&questions.QuestionIn{
			Content: "this is another question?",
			Round:   "pair",
		}, http.StatusCreated,
	},
	{
		"Add a question to quibly and to round answer, language de",
		"quibly",
		&questions.QuestionIn{
			Content:      "what is the funniest thing ever told?",
			LanguageCode: "de",
			Round:        "answers",
		}, http.StatusCreated,
	},
	{
		"Add a question to quibly and to round group",
		"quibly",
		&questions.QuestionIn{
			Content: "what does ATGM stand for?",
			Round:   "group",
		}, http.StatusCreated,
	},
	{
		"Add a question to drawlosseum, language ur",
		"drawlosseum",
		&questions.QuestionIn{
			Content:      "camel",
			LanguageCode: "ur",
		}, http.StatusCreated,
	},
	{
		"Add another question to drawlosseum",
		"drawlosseum",
		&questions.QuestionIn{
			Content: "pencil",
		}, http.StatusCreated,
	},
	{
		"Add yet another question to drawlosseum",
		"drawlosseum",
		&questions.QuestionIn{
			Content: "food fight",
		}, http.StatusCreated,
	},
	{
		"Add a question to fibbing it, round opinion new group bike group, language en",
		"fibbing_it",
		&questions.QuestionIn{
			Content:      "do you love bikes?",
			LanguageCode: "en",
			Round:        "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "question",
			},
		}, http.StatusCreated,
	},
	{
		"Add another question to fibbing it, round opinion new group bike group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "how much does liam love bikes?",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "question",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round opinion existing group bike group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "super love",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "answer",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round free_form existing group bike group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What is love?",
			Round:   "free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round free_form new group horse group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What is the fastest horse?",
			Round:   "free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "horse_group",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round free_form existing group horse group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What is the second horse called?",
			Round:   "free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "horse_group",
			},
		}, http.StatusCreated,
	},
	{
		"Add an answer to fibbing it, round likely",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "to never eat a vegetable again?",
			Round:   "likely",
		}, http.StatusCreated,
	},
	{
		"Add question to quibly, invalid round",
		"quibly",
		&questions.QuestionIn{
			Content: "This is another question?",
			Round:   "invalid",
		}, http.StatusBadRequest,
	},
	{
		"Add question to quibly, invalid2 round",
		"quibly",
		&questions.QuestionIn{
			Content: "This is another question?",
			Round:   "invalid2",
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, invalid round free_form",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What is the fastest horse?",
			Round:   "invalid_free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "horse_group",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, invalid language code",
		"fibbing_it",
		&questions.QuestionIn{
			Content:      "What is the fastest horse?",
			LanguageCode: "deed",
			Round:        "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "horse_group",
				Type: "answer",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, round opinion invalid answers type",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "super love",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "answers",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, round opinion invalid questions type",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "super love",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "questions",
			},
		}, http.StatusBadRequest,
	},
	{
		"Add an answer to fibbing it, round opinion invalid type",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "super love",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "type",
			},
		}, http.StatusBadRequest,
	},
	{
		"game does not exist but missing content",
		"quibly v3",
		&questions.QuestionIn{}, http.StatusBadRequest,
	},
	{
		"game does not exist",
		"quibly_v2",
		&questions.QuestionIn{
			Content: "super love",
		}, http.StatusNotFound,
	},
	{
		"another game does not exist",
		"quibly v3",
		&questions.QuestionIn{
			Content: "super love",
		}, http.StatusNotFound,
	},
	{
		"Add a question that already exists to quibly and to round pair",
		"quibly",
		&questions.QuestionIn{
			Content: "this is also question?",
			Round:   "pair",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to quibly and to round answer",
		"quibly",
		&questions.QuestionIn{
			Content: "pink mustard",
			Round:   "answers",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to quibly and to round answer",
		"quibly",
		&questions.QuestionIn{
			Content:      "german",
			LanguageCode: "de",
			Round:        "answers",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to quibly and to round group",
		"quibly",
		&questions.QuestionIn{
			Content: "what does ATGM stand for?",
			Round:   "group",
		}, http.StatusConflict,
	},
	{
		"Add a question that already exists to drawlosseum",
		"drawlosseum",
		&questions.QuestionIn{
			Content: "horse",
		}, http.StatusConflict,
	},
	{
		"Add another question that already exists to drawlosseum",
		"drawlosseum",
		&questions.QuestionIn{
			Content: "pencil",
		}, http.StatusConflict,
	},
	{
		"Add yet another question that already exists to drawlosseum",
		"drawlosseum",
		&questions.QuestionIn{
			Content: "food fight",
		}, http.StatusConflict,
	},
	{
		"Add a question to fibbing it that already exists, round opinion new group bike group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "do you love bikes?",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "question",
			},
		}, http.StatusConflict,
	},
	{
		"Add another question to fibbing it that already exists, round opinion new group bike group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "how much does liam love bikes?",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "question",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round opinion existing group bike group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "super love",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
				Type: "answer",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round free_form existing group bike group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What is love?",
			Round:   "free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "bike_group",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round free_form new group horse group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What is the fastest horse?",
			Round:   "free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "horse_group",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it that already exists, round free_form existing group horse group",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What is the second horse called?",
			Round:   "free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "horse_group",
			},
		}, http.StatusConflict,
	},
	{
		"Add an answer to fibbing it tthat already exists, round likely",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "to never eat a vegetable again?",
			Round:   "likely",
		}, http.StatusConflict,
	},
	{
		"Add a question to fibbing it that already exists",
		"fibbing_it",
		&questions.QuestionIn{
			Content: "What do you think about horses?",
			Round:   "opinion",
			Group: &questions.QuestionGroupInOut{
				Name: "horse_group",
				Type: "question",
			},
		}, http.StatusConflict,
	},
}

var RemoveQuestion = []struct {
	TestDescription string
	Game            string
	ID              string
	Expected        int
}{
	{
		"Remove a question from quibly and round pair",
		"quibly",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		http.StatusOK,
	},
	{
		"Remove a question from drawlossuem",
		"drawlosseum",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		http.StatusOK,
	},
	{
		"Remove a question from fibbing it",
		"fibbing_it",
		"714464a5-337f-4ce7-a4df-2b00764e5c5b",
		http.StatusOK,
	},
	{
		"Remove a question from fibbing it that was already removed",
		"fibbing_it",
		"714464a5-337f-4ce7-a4df-2b00764e5c5b",
		http.StatusNotFound,
	},
	{
		"Remove a question from quibly that was already removed",
		"quibly",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		http.StatusNotFound,
	},
	{
		"Remove a question that was already removed from drawlossuem",
		"drawlosseum",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		http.StatusNotFound,
	},
	{
		"Remove a question that doesn't exist from fibbing_it",
		"fibbing_it",
		"invalid-id",
		http.StatusNotFound,
	},
}

var AddTranslationQuestion = []struct {
	TestDescription string
	Game            string
	LanguageCode    string
	ID              string
	Payload         interface{}
	Expected        int
}{
	{
		"Update question in quibly and round pair, new language fr",
		"quibly",
		"fr",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		&questions.QuestionTranslationIn{
			Content: "this is a question?",
		},
		http.StatusCreated,
	},
	{
		"Update question in quibly and round pair, replace existing language de",
		"quibly",
		"de",
		"bf64d60c-62ee-420a-976e-bfcaec77ad8b",
		&questions.QuestionTranslationIn{
			Content: "le german?",
		},
		http.StatusCreated,
	},
	{
		"Update question in drawlosseum",
		"drawlosseum",
		"hi",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		&questions.QuestionTranslationIn{
			Content: "ऊंट",
		},
		http.StatusCreated,
	},
	{
		"Update question in fibbing it, round opinion",
		"fibbing_it",
		"it",
		"580aeb14-d907-4a22-82c8-f2ac544a2cd1",
		&questions.QuestionTranslationIn{
			Content: "Cosa ne pensi dei cavalli?",
		},
		http.StatusCreated,
	},
	{
		"Update question in fibbing it, round opinion and answers section",
		"fibbing_it",
		"de",
		"aa9fe2b5-79b5-458d-814b-45ff95a617fc",
		&questions.QuestionTranslationIn{
			Content: "Liebe",
		}, http.StatusCreated,
	},
	{
		"Missing content",
		"quibly",
		"en",
		"a9c00e19-d41e-4b15-a8bd-ec921af9123d",
		&questions.QuestionIn{}, http.StatusBadRequest,
	},
	{
		"Update question in fibbing it but invalid language code",
		"fibbing_it",
		"ittt",
		"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
		&questions.QuestionTranslationIn{
			Content: "was ist Liebe?",
		}, http.StatusBadRequest,
	},
	{
		"game does not exist",
		"quibly v3",
		"de",
		"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
		&questions.QuestionTranslationIn{
			Content: "was ist Liebe?",
		}, http.StatusNotFound,
	},
	{
		"Question doesn't exist",
		"fibbing_it",
		"de",
		"9f64d60c-62ee-420a-976e-bfcaec77ad8b",
		&questions.QuestionTranslationIn{
			Content: "was ist Liebe?",
		}, http.StatusNotFound,
	},
}

var RemoveTranslationQuestion = []struct {
	TestDescription string
	Game            string
	ID              string
	LanguageCode    string
	Expected        int
}{
	{
		"Delete a question quibly from round pair",
		"quibly",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		"en",
		http.StatusOK,
	},
	{
		"Delete a question quibly from round pair, language ur",
		"quibly",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		"ur",
		http.StatusOK,
	},
	{
		"Delete a question quibly from round answers",
		"quibly",
		"bf64d60c-62ee-420a-976e-bfcaec77ad8b",
		"en",
		http.StatusOK,
	},
	{
		"Delete a question quibly from round group, language fr",
		"quibly",
		"4b4dd325-04fd-4aa4-9382-2874dcfd5cae",
		"fr",
		http.StatusOK,
	},
	{
		"Delete a question drawlosseum",
		"drawlosseum",
		"815464a5-337f-4ce7-a4df-2b00764e5c6c",
		"en",
		http.StatusOK,
	},
	{
		"Delete another question drawlosseum",
		"drawlosseum",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		"en",
		http.StatusOK,
	},
	{
		"Delete a question to fibbing it, round opinion from group horse group",
		"fibbing_it",
		"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
		"en",
		http.StatusOK,
	},
	{
		"Delete a answer to fibbing it, round opinion from group horse group",
		"fibbing_it",
		"03a462ba-f483-4726-aeaf-b8b6b03ce3e2",
		"en",
		http.StatusOK,
	},
	{
		"Delete a question quibly from round pair that was already deleted",
		"quibly",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		"en",
		http.StatusNotFound,
	},
	{
		"Delete a question drawlosseum that was already deleted",
		"drawlosseum",
		"815464a5-337f-4ce7-a4df-2b00764e5c6c",
		"en",
		http.StatusNotFound,
	},
	{
		"Delete a question already removed from fibbing it, round free_form from group bike group",
		"fibbing_it",
		"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
		"en",
		http.StatusNotFound,
	},
	{
		"Delete another  already removed from fibbing it, round likely",
		"fibbing_it",
		"03a462ba-f483-4726-aeaf-b8b6b03ce3e2",
		"en",
		http.StatusNotFound,
	},
}

var EnableQuestion = []struct {
	TestDescription string
	Game            string
	ID              string
	Expected        int
}{
	{
		"Enable a question, quibly and round pair",
		"quibly",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		http.StatusOK,
	},
	{
		"Enable a question, quibly and round answers",
		"quibly",
		"4b4dd325-04fd-4aa4-9382-2874dcfd5cae",
		http.StatusOK,
	},
	{
		"Enable a question, fibbing_it and round opinion",
		"fibbing_it",
		"7799e38a-758d-4a1b-a191-99c59440af76",
		http.StatusOK,
	},
	{
		"Enable an answer, fibbing_it and round opinion",
		"fibbing_it",
		"03a462ba-f483-4726-aeaf-b8b6b03ce3e2",
		http.StatusOK,
	},
	{
		"Enable a question, fibbing_it and round free_form",
		"fibbing_it",
		"580aeb14-d907-4a22-82c8-f2ac544a2cd1",
		http.StatusOK,
	},
	{
		"Enable a question, fibbing_it and round likely",
		"fibbing_it",
		"d80f2d90-0fb0-462a-8fbd-1aa00b4e42a5",
		http.StatusOK,
	},
	{
		"Enable a question, drawlosseum",
		"drawlosseum",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		http.StatusOK,
	},
	{
		"Enable an already enabled question, drawlosseum",
		"drawlosseum",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		http.StatusOK,
	},
	{
		"Game does not exist",
		"quibly v3",
		"901464a5-337f-4ce7-a4df-2b00764e5d8d",
		http.StatusNotFound,
	},
}

var DisableQuestion = []struct {
	TestDescription string
	Game            string
	ID              string
	Expected        int
}{
	{
		"Disable a question, quibly and round pair",
		"quibly",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		http.StatusOK,
	},
	{
		"Disable a question, quibly and round answers",
		"quibly",
		"bf64d60c-62ee-420a-976e-bfcaec77ad8b",
		http.StatusOK,
	},
	{
		"Disable a question, fibbing_it and round opinion",
		"fibbing_it",
		"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
		http.StatusOK,
	},
	{
		"Disable an answer, fibbing_it and round opinion",
		"fibbing_it",
		"03a462ba-f483-4726-aeaf-b8b6b03ce3e2",
		http.StatusOK,
	},
	{
		"Disable a question, fibbing_it and round likely",
		"fibbing_it",
		"d5aa9153-f48c-45cc-b411-fb9b2d38e78f",
		http.StatusOK,
	},
	{
		"Disable a question, drawlosseum",
		"drawlosseum",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		http.StatusOK,
	},
	{
		"Question does not exist",
		"quibly",
		"90aa9153-f48c-45cc-b411-fb9b2d38e78f",
		http.StatusNotFound,
	},
	{
		"Game does not exist",
		"quibly v3",
		"d5aa9153-f48c-45cc-b411-fb9b2d38e78f",
		http.StatusNotFound,
	},
}

var GetAllGroups = []struct {
	TestDescription string
	Payload         *questions.GroupInput
	ExpectedGroups  []string
	ExpectedCode    int
}{
	{
		"Get all groups from questions from the opinion round in the Fibbing It game",
		&questions.GroupInput{
			GameParams: internal.GameParams{
				Name: "fibbing_it",
			},
			RoundParams: internal.RoundParams{Round: "opinion"},
		},
		[]string{
			"horse_group",
		},
		http.StatusOK,
	},

	{
		"Get all groups from questions from the free form round in the Fibbing It game",
		&questions.GroupInput{
			GameParams: internal.GameParams{
				Name: "fibbing_it",
			},
			RoundParams: internal.RoundParams{Round: "free_form"},
		},
		[]string{
			"bike_group",
			"cat_group",
		},
		http.StatusOK,
	},

	{
		"Try to get groups from a round in Fibbing It that does not have groups",
		&questions.GroupInput{
			GameParams: internal.GameParams{
				Name: "fibbing_it",
			},
			RoundParams: internal.RoundParams{Round: "likely"},
		},
		[]string{},
		http.StatusNotFound,
	},

	{
		"Try to get groups from a non-existent round",
		&questions.GroupInput{
			GameParams: internal.GameParams{
				Name: "fibbing_it",
			},
			RoundParams: internal.RoundParams{Round: "genocide"},
		},
		[]string{},
		http.StatusNotFound,
	},

	{
		"Try to get groups from a game that does not have groups",
		&questions.GroupInput{
			GameParams: internal.GameParams{
				Name: "quibly",
			},
			RoundParams: internal.RoundParams{Round: "opinion"},
		},
		[]string{},
		http.StatusNotFound,
	},
}

var GetQuestions = []struct {
	TestDescription   string
	Game              string
	Round             string
	Language          string
	Limit             int
	GroupName         string
	Enabled           string
	Random            bool
	ExpectedStatus    int
	ExpectedQuestions []questions.QuestionOut
}{
	{
		"Get some quibly questions for round pair",
		"quibly",
		"pair",
		"",
		2,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "this is a question?",
				Type:    "question",
			},
			{
				Content: "this is also question?",
				Type:    "question",
			},
		},
	},
	{
		"Get some quibly questions for round answers",
		"quibly",
		"answers",
		"de",
		1,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "german",
				Type:    "answer",
			},
		},
	},
	{
		"Get some quibly questions for round group",
		"quibly",
		"group",
		"fr",
		1,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "this is a another question?",
				Type:    "question",
			},
		},
	},
	{
		"Get some quibly questions for round group and limit 10",
		"quibly",
		"group",
		"fr",
		10,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "this is a another question?",
				Type:    "question",
			},
		},
	},
	{
		"Get some quibly questions for round group and no language for that question",
		"quibly",
		"group",
		"ur",
		10,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{},
	},
	{
		"Get some fibbing_it questions for round opinion",
		"fibbing_it",
		"opinion",
		"en",
		3,
		"horse_group",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "What do you think about horses?",
				Type:    "question",
			},
			{
				Content: "What do you think about camels?",
				Type:    "question",
			},
			{
				Content: "cool",
				Type:    "answer",
			},
			{
				Content: "tasty",
				Type:    "answer",
			},
			{
				Content: "lame",
				Type:    "answer",
			},
		},
	},
	{
		"Get some fibbing_it questions for round free_form (no language specified, defaults to en)",
		"fibbing_it",
		"free_form",
		"",
		5,
		"bike_group",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "Favourite bike colour?",
				Type:    "question",
			},
			{
				Content: "A funny question?",
				Type:    "question",
			},
		},
	},
	{
		"Get some fibbing_it questions for round likely",
		"fibbing_it",
		"likely",
		"",
		5,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "to eat ice-cream from the tub",
				Type:    "question",
			},
			{
				Content: "to get arrested",
				Type:    "question",
			},
		},
	},
	{
		"Get some drawlosseum questions",
		"drawlosseum",
		"drawing",
		"",
		2,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "horse",
				Type:    "answer",
			},
			{
				Content: "spoon",
				Type:    "answer",
			},
		},
	},
	{
		"Get some drawlosseum questions, bad limit -1",
		"drawlosseum",
		"drawing",
		"",
		-1,
		"",
		"",
		false,
		http.StatusBadRequest,
		[]questions.QuestionOut{},
	},
	{
		"Get some drawlosseum round doesn't exist",
		"drawlosseum",
		"drawing3",
		"",
		2,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{},
	},
	{
		"Get some drawlosseum round doesn't exist",
		"drawlosseum",
		"drawing3",
		"",
		2,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{},
	},
	{
		"Get some questions game doesn't exist",
		"game",
		"abc",
		"",
		2,
		"",
		"",
		false,
		http.StatusOK,
		[]questions.QuestionOut{},
	},
	{
		"Get some fibbing_it questions for round opinion, enabled",
		"fibbing_it",
		"opinion",
		"en",
		3,
		"horse_group",
		"enabled",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "What do you think about horses?",
				Type:    "question",
			},
			{
				Content: "What do you think about camels?",
				Type:    "question",
			},
			{
				Content: "cool",
				Type:    "answer",
			},
			{
				Content: "tasty",
				Type:    "answer",
			},
			{
				Content: "lame",
				Type:    "answer",
			},
		},
	},
	{
		"Get some fibbing_it questions for round opinion, not enabled",
		"fibbing_it",
		"opinion",
		"en",
		3,
		"horse_group",
		"disabled",
		false,
		http.StatusOK,
		[]questions.QuestionOut{},
	},
	{
		"Get some fibbing_it disabled questions",
		"fibbing_it",
		"free_form",
		"it",
		3,
		"cat_group",
		"disabled",
		false,
		http.StatusOK,
		[]questions.QuestionOut{
			{
				Content: "Perché sono superiori i gatti di Liam?",
				Type:    "question",
			},
		},
	},
	{
		"Get some quibly questions for round pair and random",
		"quibly",
		"pair",
		"",
		2,
		"",
		"",
		true,
		http.StatusOK,
		[]questions.QuestionOut{},
	},
}

var GetQuestionById = []struct {
	TestDescription string
	Game            string
	LanguageCode    string
	ID              string
	ExpectedPayload questions.QuestionGenericOut
	Expected        int
}{
	{
		"Get a quibly question",
		"quibly",
		"de",
		"4d18ac45-8034-4f8e-b636-cf730b17e51a",
		questions.QuestionGenericOut{
			Round:   "pair",
			Enabled: true,
			Content: "this is a question?",
		},
		http.StatusOK,
	},
	{
		"Get a drawlosseum question",
		"drawlosseum",
		"en",
		"101464a5-337f-4ce7-a4df-2b00764e5d8d",
		questions.QuestionGenericOut{
			Round:   "drawing",
			Enabled: true,
			Content: "spoon",
		},
		http.StatusOK,
	},
	{
		"Get a fibbing it question",
		"fibbing_it",
		"it",
		"d80f2d90-0fb0-462a-8fbd-1aa00b4e42a5",
		questions.QuestionGenericOut{
			Content: "Perché sono superiori i gatti di Liam?",
			Enabled: false,
			Round:   "free_form",
			Group: &questions.QuestionGroupInOut{
				Name: "cat_group",
			},
		},
		http.StatusOK,
	},
	{
		"Game does not exist",
		"quibly v3",
		"de",
		"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
		questions.QuestionGenericOut{}, http.StatusNotFound,
	},
	{
		"Question does not exist",
		"fibbing_it",
		"de",
		"9f64d60c-62ee-420a-976e-bfcaec77ad8b",
		questions.QuestionGenericOut{}, http.StatusNotFound,
	},
	{
		"Question does not exist",
		"fibbing_it",
		"de",
		"1010d60c-62ee-420a-976e-bfcaec771010",
		questions.QuestionGenericOut{}, http.StatusNotFound,
	},
}

var GetAllQuestionsIds = []struct {
	TestDescription string
	Game            string
	Limit           int64
	Cursor          string
	ExpectedPayload questions.AllQuestionOut
	ExpectedStatus  int
}{
	{
		"Get all questions from fibbing it",
		"fibbing_it",
		5,
		"",
		questions.AllQuestionOut{
			IDs: []string{
				"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
				"7799e38a-758d-4a1b-a191-99c59440af76",
				"03a462ba-f483-4726-aeaf-b8b6b03ce3e2",
				"d5aa9153-f48c-45cc-b411-fb9b2d38e78f",
				"138bc208-2849-41f3-bbd8-3226a96c5370",
			},
			Cursor: "138bc208-2849-41f3-bbd8-3226a96c5370",
		},
		http.StatusOK,
	},
	{
		"Get all questions from fibbing it using pagination",
		"fibbing_it",
		5,
		"138bc208-2849-41f3-bbd8-3226a96c5370",
		questions.AllQuestionOut{
			IDs: []string{
				"3e2889f6-56aa-4422-a7c5-033eafa9fd39",
				"7799e38a-758d-4a1b-a191-99c59440af76",
				"d5aa9153-f48c-45cc-b411-fb9b2d38e78f",
				"580aeb14-d907-4a22-82c8-f2ac544a2cd1",
				"aa9fe2b5-79b5-458d-814b-45ff95a617fc",
			},
			Cursor: "aa9fe2b5-79b5-458d-814b-45ff95a617fc",
		},
		http.StatusOK,
	},
	{
		"Get all questions from drawlossuem it using pagination",
		"drawlosseum",
		5,
		"",
		questions.AllQuestionOut{
			IDs: []string{
				"815464a5-337f-4ce7-a4df-2b00764e5c6c",
				"101464a5-337f-4ce7-a4df-2b00764e5d8d",
			},
			Cursor: "",
		},
		http.StatusOK,
	},
	{
		"Get all questions from quibly",
		"quibly",
		2,
		"",
		questions.AllQuestionOut{
			IDs: []string{
				"4d18ac45-8034-4f8e-b636-cf730b17e51a",
				"a9c00e19-d41e-4b15-a8bd-ec921af9123d",
			},
			Cursor: "a9c00e19-d41e-4b15-a8bd-ec921af9123d",
		},
		http.StatusOK,
	},
	{
		"Get all questions from quibly it using pagination",
		"quibly",
		3,
		"a9c00e19-d41e-4b15-a8bd-ec921af9123d",
		questions.AllQuestionOut{
			IDs: []string{
				"bf64d60c-62ee-420a-976e-bfcaec77ad8b",
			},
			Cursor: "",
		},
		http.StatusOK,
	},
	{
		"Get with invalid limit",
		"quibly",
		-1,
		"",
		questions.AllQuestionOut{
			IDs:    []string{},
			Cursor: "",
		},
		http.StatusBadRequest,
	},
}
