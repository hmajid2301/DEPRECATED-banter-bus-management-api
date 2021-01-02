package biz

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
)

// GameService is struct data required by all game service functions.
type GameService struct {
	DB core.Repository
}

// Add is used to add a new game.
func (game *GameService) Add(name string, rulesURL string) error {
	exists := game.doesItExist(name)
	if exists {
		return errors.AlreadyExistsf("The game %s", name)
	}

	gameInfo, err := getGameType(name, models.GenericQuestion{}, game.DB)
	if err != nil {
		return err
	}

	inserted, err := gameInfo.AddGame(rulesURL)
	if !inserted {
		return errors.Errorf("Failed to add the new game %s", name)
	}
	return err
}

// Get is used to get specific information about a game.
func (game *GameService) Get(name string) (*models.GameInfo, error) {
	var (
		filter   = &models.GameInfo{Name: name}
		gameInfo *models.GameInfo
	)

	err := game.DB.Get("game", filter, &gameInfo)
	if err != nil {
		return &models.GameInfo{}, err
	}

	return gameInfo, nil
}

// GetAll is used to get all game.
func (game *GameService) GetAll(enabled *bool) ([]string, error) {
	var games []*models.GameInfo

	err := game.DB.GetAll("game", &games)
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

// Remove is used to remove a game.
func (game *GameService) Remove(name string) error {
	exists := game.doesItExist(name)
	if !exists {
		return errors.NotFoundf("The game %s", name)
	}

	filter := &models.GameInfo{Name: name}
	deleted, err := game.DB.Delete("game", filter)
	if !deleted || err != nil {
		return errors.Errorf("Failed to remove game %s", name)
	}

	return nil
}

// UpdateEnable is used to update the enable state of a game.
func (game *GameService) UpdateEnable(name string, enabled bool) (bool, error) {
	filter := &models.GameInfo{Name: name}
	gameInfo, err := game.Get(name)

	if gameInfo.Name == "" || err != nil {
		return false, errors.NotFoundf("The game %s", name)
	}

	update := &models.GameInfo{Name: name, RulesURL: gameInfo.RulesURL, Enabled: &enabled}
	updated, err := game.DB.UpdateEntry("game", filter, update)
	if err != nil {
		return false, errors.Errorf("Failed to update game %s", err)
	}
	return updated, err
}

func (game *GameService) doesItExist(name string) bool {
	gameInfo, err := game.Get(name)
	if err != nil {
		return false
	}

	return gameInfo.Name != ""
}

func getGameType(name string, question models.GenericQuestion, db core.Repository) (models.PlayableGame, error) {
	var playableGame models.PlayableGame
	switch name {
	case "quibly":
		playableGame = games.Quibly{CurrentQuestion: question, DB: db}
	case "fibbing_it":
		playableGame = games.FibbingIt{CurrentQuestion: question, DB: db}
	case "drawlosseum":
		playableGame = games.Drawlosseum{DB: db}
	default:
		return nil, errors.BadRequestf("Invalid game %s", name)
	}

	return playableGame, nil
}
