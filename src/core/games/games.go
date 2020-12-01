package games

// PlayableGame interface for all game types.
type PlayableGame interface {
	AddGame(string) (bool, error)
	GetQuestionPath() string
	ValidateQuestionInput() error
}
