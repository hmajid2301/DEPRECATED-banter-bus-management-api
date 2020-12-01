package models

// User struct to hold a database entry for a user, i.e a player with an account
type User struct {
	Username      string           `bson:"username"`
	Admin         *bool            `bson:"admin,omitempty"`
	Privacy       string           `bson:"privacy,omitempty"`
	Membership    string           `bson:"membership,omitempty"`
	Preferences   *UserPreferences `bson:"preferences,omitempty"`
	Friends       []Friend         `bson:"friends,omitempty"`
	Stories       []Story          `bson:"stories,omitempty"`
	QuestionPools []QuestionPool   `bson:"question_pools,omitempty"`
}

// UserPreferences struct to hold information on a user's preferences
type UserPreferences struct {
	LanguageCode string `bson:"language_code,omitempty"`
}

// Friend struct to hold details about a user's friend
type Friend struct {
	Username string `bson:"username"`
}
