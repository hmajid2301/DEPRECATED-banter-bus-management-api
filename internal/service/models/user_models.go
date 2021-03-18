package models

import "gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"

// User struct to hold a database entry for a user, i.e a player with an account
type User struct {
	Username      string           `bson:"username"`
	Admin         *bool            `bson:"admin,omitempty"`
	Privacy       string           `bson:"privacy,omitempty"`
	Membership    string           `bson:"membership,omitempty"`
	Preferences   *UserPreferences `bson:"preferences,omitempty"`
	Friends       []Friend         `bson:"friends"`
	QuestionPools []Pool           `bson:"question_pools"        json:"question_pools"`
}

// Add is used to add users to a database.
func (user *User) Add(db database.Database) (bool, error) {
	inserted, err := db.Insert("user", user)
	return inserted, err
}

// Get is used to retrieve a user from the database.
func (user *User) Get(db database.Database, filter map[string]string) error {
	err := db.Get("user", filter, user)
	return err
}

// Update is used to update a user in the database.
func (user *User) Update(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.Update("user", filter, user)
	return updated, err
}

// Users is a list of users.
type Users []User

// Add method adds (a list of) users at once.
func (users *Users) Add(db database.Database) error {
	err := db.InsertMultiple("user", users)
	return err
}

// Get method gets all the users in the user collection.
func (users *Users) Get(db database.Database, filter map[string]string) error {
	err := db.GetAll("user", filter, users)
	return err
}

// Delete is used to delete a list of users that match a filter.
func (users Users) Delete(db database.Database, filter map[string]string) (bool, error) {
	deleted, err := db.DeleteAll("user", filter)
	return deleted, err
}

// ToInterface converts users (list of users) into a list of interfaces, required by GetAll MongoDB.
func (users Users) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(users))
	for i, item := range users {
		interfaceObject[i] = item
	}
	return interfaceObject
}

// UserPreferences struct to hold information on a user's preferences
type UserPreferences struct {
	LanguageCode string `bson:"language_code,omitempty" json:"language_code"`
}

// Friend struct to hold details about a user's friend
type Friend struct {
	Username string `bson:"username"`
}
