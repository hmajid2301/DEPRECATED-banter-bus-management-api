package data

import (
	"net/http"

	serverModels "banter-bus-server/src/server/models"
)

var admin = true
var notAdmin = !admin

//AddUser is the test data for adding new users.
var AddUser = []struct {
	TestDescription string
	Payload         interface{}
	ExpectedStatus  int
}{
	{
		"Add a new user",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "paid",
			Admin:      &admin,
		},
		http.StatusCreated,
	},

	{
		"Add a user that has already been added",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "paid",
			Admin:      &admin,
		},
		http.StatusConflict,
	},

	{
		"Try to add a user with a taken username but different admin privileges",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "paid",
			Admin:      &notAdmin,
		},
		http.StatusConflict,
	},

	{
		"Try to add a user with a taken username but different membership",
		&serverModels.NewUser{
			Username:   "lmoz25",
			Membership: "free",
			Admin:      &admin,
		},
		http.StatusConflict,
	},

	{
		"Try to add a user with an invalid membership",
		&serverModels.NewUser{
			Username:   "seeb123",
			Membership: "overlord",
			Admin:      &notAdmin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with no username",
		&serverModels.NewUser{
			Username:   "",
			Membership: "paid",
			Admin:      &admin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with invalid fields 1",
		struct {
			Usernam, Membership string
			Admin               *bool
		}{
			Usernam:    "seeb123",
			Membership: "free",
			Admin:      &notAdmin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with invalid fields 2",
		struct {
			Username, Membershap string
			Admin                *bool
		}{
			Username:   "seeb123",
			Membershap: "free",
			Admin:      &notAdmin,
		},
		http.StatusBadRequest,
	},

	{
		"Try to add a user with invalid fields 3",
		struct {
			Username, Membership string
			Admine               *bool
		}{
			Username:   "seeb123",
			Membership: "free",
			Admine:     &notAdmin,
		},
		http.StatusBadRequest,
	},
}

//GetAllUsers is the test data for geting all users
var GetAllUsers = []struct {
	TestDescription   string
	Filter            *serverModels.ListUserParams
	ExpectedUsernames []string
}{
	{
		"Get all users no filter",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "all",
			Membership:  "all",
		},
		[]string{
			"virat_kohli",
			"roh1t_sharma",
			"dhawanShikhar",
		},
	},
	{
		"Get all admins",
		&serverModels.ListUserParams{
			AdminStatus: "admin",
			Privacy:     "all",
			Membership:  "all",
		},
		[]string{
			"roh1t_sharma",
		},
	},
	{
		"Get all non-admins",
		&serverModels.ListUserParams{
			AdminStatus: "non_admin",
			Privacy:     "all",
			Membership:  "all",
		},
		[]string{
			"virat_kohli",
			"dhawanShikhar",
		},
	},
	{
		"Get all members with public profiles",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "public",
			Membership:  "all",
		},
		[]string{
			"virat_kohli",
		},
	},
	{
		"Get all members with private profiles",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "private",
			Membership:  "all",
		},
		[]string{
			"roh1t_sharma",
			"dhawanShikhar",
		},
	},
	{
		"Get all members with free membership",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "all",
			Membership:  "free",
		},
		[]string{
			"roh1t_sharma",
		},
	},
	{
		"Get all members with paid membership",
		&serverModels.ListUserParams{
			AdminStatus: "all",
			Privacy:     "all",
			Membership:  "paid",
		},
		[]string{
			"virat_kohli",
			"dhawanShikhar",
		},
	},
	{
		"Get members who are admin, have a private account and free membership",
		&serverModels.ListUserParams{
			AdminStatus: "admin",
			Privacy:     "private",
			Membership:  "free",
		},
		[]string{
			"roh1t_sharma",
		},
	},
}

//GetUser is the test data for getting users.
var GetUser = []struct {
	TestDescription string
	Username        string
	ExpectedStatus  int
	ExpectedUser    serverModels.User
}{
	{
		"Get a user",
		"virat_kohli",
		http.StatusOK,
		serverModels.User{
			Username:   "virat_kohli",
			Admin:      &notAdmin,
			Privacy:    "public",
			Membership: "paid",
			Preferences: &serverModels.UserPreferences{
				LanguageCode: "pa",
			},
			Friends: []serverModels.Friend{
				{
					Username: "roh1t_sharma",
				},
				{
					Username: "dhawanShikhar",
				},
			},
			Stories: []serverModels.Story{
				{
					GameName: "quibly",
					Question: "how many fish are there?",
					Answers: []serverModels.StoryAnswer{
						{
							Answer: "one",
							Votes:  12341,
						},
						{
							Answer: "many",
							Votes:  0,
						},
					},
				},
			},
		},
	},
	{
		"Get another user",
		"roh1t_sharma",
		http.StatusOK,
		serverModels.User{
			Username:   "roh1t_sharma",
			Admin:      &admin,
			Privacy:    "private",
			Membership: "free",
			Preferences: &serverModels.UserPreferences{
				LanguageCode: "mr",
			},
			Friends: []serverModels.Friend{
				{
					Username: "virat_kohli",
				},
			},
		},
	},
	{
		"Get a user that doesn't exists",
		"azharAli",
		http.StatusNotFound,
		serverModels.User{},
	},
}
