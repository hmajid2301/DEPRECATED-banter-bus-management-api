package serverModels

// NewUser struct is the data needed to add a new user via the API
type NewUser struct {
	Username   string `json:"username"   description:"The screen name of the user." example:"lmoz25" validate:"required"`
	Membership string `json:"membership" description:"membership for a user"        example:"free"   validate:"required,oneof=free paid"`
	Admin      *bool  `json:"admin"      description:"Whether the user is admin"                     validate:"required"`
}

// User struct to hold a database entry for a user, i.e a player with an account
type User struct {
	Username    string           `json:"username"    description:"The screen name of the user."          example:"lmoz25"`
	Admin       *bool            `json:"admin"       description:"Whether or not the user is admin."`
	Privacy     string           `json:"privacy"     description:"The privacy level of the user."        example:"public"`
	Membership  string           `json:"membership"  description:"The level of membership the user has." example:"free"`
	Preferences *UserPreferences `json:"preferences" description:"Collection of user preferences."`
	Friends     []Friend         `json:"friends"     description:"List of friends the user has."`
}

// UserPreferences struct to hold information on a user's preferences
type UserPreferences struct {
	LanguageCode string `json:"language_code" description:"Details of the user's preferred display and question language stored using ISO 2-letter language abbreviations" example:"en"`
}

// Friend struct to hold details about a user's friend
type Friend struct {
	Username string `json:"username" description:"The screen name of the friend" example:"seeb123"`
}

// NewQuestionPool is the data required to create a new question pool.
type NewQuestionPool struct {
	PoolName     string `json:"pool_name"     description:"The unique name of the question pool."                                                                          example:"my_pool" validate:"required"`
	GameName     string `json:"game_name"     description:"The type of game the question pool pertains to"                                                                 example:"quibly"  validate:"required,oneof=quibly fibbing_it drawlosseum" enum:"quibly,fibbing_it,drawlosseum"`
	LanguageCode string `json:"language_code" description:"Details of the user's preferred display and question language stored using ISO 2-letter language abbreviations" example:"en"                                                                                                   default:"en"`
	Privacy      string `json:"privacy"       description:"The privacy setting for this question pool"                                                                                       validate:"required,oneof=public private friends"        enum:"public,private,friends"`
}

// UpdateQuestionPool is the http body data required to update a question i.e. add/remove a question from the pool.
type UpdateQuestionPool struct {
	NewQuestion
	Operation string `json:"operation" description:"The update operation type." validate:"required,oneof=add remove" example:"add" enum:"add,remove"`
}

// QuestionPool is the list of custom questions for different game.
type QuestionPool struct {
	PoolName     string                `json:"pool_name"     description:"The unique name of the question pool."                                                                          example:"my_pool"`
	GameName     string                `json:"game_name"     description:"The type of game the story pertains to"                                                                         example:"quibly"`
	LanguageCode string                `json:"language_code" description:"Details of the user's preferred display and question language stored using ISO 2-letter language abbreviations" example:"en"`
	Privacy      string                `json:"privacy"       description:"The privacy setting for this question pool"                                                                                       enum:"public,private,friends"`
	Questions    QuestionPoolQuestions `json:"questions"     description:"List of questions in this pool."`
}

// QuestionPoolQuestions contains the user's questions (only one game).
type QuestionPoolQuestions struct {
	Drawlosseum DrawlosseumQuestionsPool `json:"drawlosseum,omitempty"`
	Quibly      QuiblyQuestionsPool      `json:"quibly,omitempty"`
	FibbingIt   FibbingItQuestionsPool   `json:"fibbing_it,omitempty"`
}

// DrawlosseumQuestionsPool is how the Drawlosseum questions are stored in a single list of possible drawing objects.
type DrawlosseumQuestionsPool struct {
	Drawings []string `json:"drawings,omitempty"`
}

// QuiblyQuestionsPool is how the Quibly questions are stored, questions are stored in different rounds of th game.
// So that the game knows which round the question belongs to, and when to "potentially" use that question.
type QuiblyQuestionsPool struct {
	Pair    []string `json:"pair,omitempty"`
	Answers []string `json:"answers,omitempty"`
	Group   []string `json:"group,omitempty"`
}

// FibbingItQuestionsPool is how the Fibbing It questions are stored, questions are stored in different rounds of the game. Note questions can be in groups.
// So that the game knows which round the question belongs to, and when to "potentially" use that question.
type FibbingItQuestionsPool struct {
	Opinion  map[string]map[string][]string `json:"opinion,omitempty"`
	FreeForm map[string][]string            `json:"free_form,omitempty"`
	Likely   []string                       `json:"likely,omitempty"`
}

// Story struct to contain information about a user story
type Story struct {
	Question string `json:"question"`
	Round    string `json:"round,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	StoryAnswers
}

// StoryAnswers contains all the different stories answers that are supported, for different game.
type StoryAnswers struct {
	Drawlosseum []StoryDrawlosseum `json:"drawlosseum,omitempty"`
	Quibly      []StoryQuibly      `json:"quibly,omitempty"`
	FibbingIt   []StoryFibbingIt   `json:"fibbing_it,omitempty"`
}

// StoryFibbingIt contains information about the Fibbing It answers for user stories
type StoryFibbingIt struct {
	Nickname string `json:"nickname,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

// StoryQuibly contains information about the Quibly answers for user stories
type StoryQuibly struct {
	Nickname string `json:"nickname,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Votes    int    `json:"votes,omitempty"`
}

// StoryDrawlosseum contains information about the Drawlosseum answers for user stories
type StoryDrawlosseum struct {
	Start DrawlosseumDrawingPoint `json:"start,omitempty"`
	End   DrawlosseumDrawingPoint `json:"end,omitempty"`
	Color string                  `json:"color,omitempty"`
}

// DrawlosseumDrawingPoint contains information about a point in a Drawlosseum drawing
type DrawlosseumDrawingPoint struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}
