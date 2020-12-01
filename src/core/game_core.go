package core

import (
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/games"
	"banter-bus-server/src/core/games/drawlosseum"
	"banter-bus-server/src/core/games/fibbingit"
	"banter-bus-server/src/core/games/quibly"
	"banter-bus-server/src/core/models"

	"github.com/juju/errors"
)

// AddGame is used add a new game.
func AddGame(name string, rulesURL string) error {
	exists := doesGameExist(name)
	if exists {
		return errors.AlreadyExistsf("The game %s", name)
	}

	game, err := getGameType(name, models.GenericQuestion{})
	if err != nil {
		return err
	}

	inserted, err := game.AddGame(rulesURL)
	if !inserted {
		return errors.Errorf("Failed to add the new game %s", name)
	}
	return err
}

// GetGame is used to get specific information about a game.
func GetGame(name string) (*models.Game, error) {
	var (
		filter = &models.Game{Name: name}
		game   *models.Game
	)

	err := database.Get("game", filter, &game)
	if err != nil {
		return &models.Game{}, err
	}

	return game, nil
}

// GetAllGames is used to get all game.
func GetAllGames(enabled *bool) ([]string, error) {
	var games []*models.Game

	err := database.GetAll("game", &games)
	if err != nil {
		return []string{}, err
	}

	var gameNames []string
	for _, game := range games {
		if enabled != nil && *enabled != *game.Enabled {
			continue
		}
		gameNames = append(gameNames, game.Name)
	}
	return gameNames, nil
}

// RemoveGame is used to remove a game.
func RemoveGame(name string) error {
	exists := doesGameExist(name)
	if !exists {
		return errors.NotFoundf("The game %s", name)
	}

	filter := &models.Game{Name: name}
	deleted, err := database.Delete("game", filter)
	if !deleted || err != nil {
		return errors.Errorf("Failed to remove game %s", name)
	}

	return nil
}

// UpdateEnableGame is used to update the enable state of a game.
func UpdateEnableGame(name string, enabled bool) (bool, error) {
	filter := &models.Game{Name: name}
	game, _ := GetGame(name)

	if game.Name == "" {
		return false, errors.NotFoundf("The game %s", name)
	}

	update := &models.Game{Name: name, RulesURL: game.RulesURL, Enabled: &enabled}
	updated, err := database.UpsertEntry("game", filter, update)
	if err != nil {
		return false, errors.Errorf("Failed to update game %s", err)
	}
	return updated, err
}

func doesGameExist(name string) bool {
	game, _ := GetGame(name)
	return game.Name != ""
}

func getGameType(name string, question models.GenericQuestion) (games.PlayableGame, error) {
	var game games.PlayableGame
	switch name {
	case "quibly":
		game = quibly.Quibly{CurrentQuestion: question}
	case "fibbing_it":
		game = fibbingit.FibbingIt{CurrentQuestion: question}
	case "drawlosseum":
		game = drawlosseum.Drawlosseum{}
	default:
		return nil, errors.BadRequestf("Invalid game: %s", name)
	}

	return game, nil
}
