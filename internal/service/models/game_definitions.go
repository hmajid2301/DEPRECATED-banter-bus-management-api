package models

// Gamer is the interface for game(s).
type Gamer interface {
	ValidateQuestion(question GenericQuestion) error
}
