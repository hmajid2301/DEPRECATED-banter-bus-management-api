package service

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// GameService is struct data required by all g service functions.
type GameService struct {
	DB   database.Database
	Name string
}

// Add is used to add a new game.
func (g *GameService) Add(rulesURL string) error {
	exists := g.doesItExist()
	if exists {
		return errors.AlreadyExistsf("The game %s", g.Name)
	}

	game, err := games.GetGame(g.Name)
	if err != nil {
		return err
	}

	newGame := game.NewGame(rulesURL)
	inserted, err := newGame.Add(g.DB)
	if !inserted {
		return errors.Errorf("Failed to add the new game %s", g.Name)
	}
	return err
}

// Get is used to get specific information about a game.
func (g *GameService) Get() (*models.Game, error) {
	var (
		filter = map[string]string{"name": g.Name}
		game   = &models.Game{}
	)

	err := game.Get(g.DB, filter)
	if err != nil {
		return &models.Game{}, err
	}

	return game, nil
}

// GetAll is used to get all games.
func (g *GameService) GetAll(enabled *bool) ([]string, error) {
	games := models.Games{}

	err := games.Get(g.DB)
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
func (g *GameService) Remove() error {
	exists := g.doesItExist()
	if !exists {
		return errors.NotFoundf("the game %s", g.Name)
	}

	filter := map[string]string{"name": g.Name}
	deleted, err := g.DB.Delete("game", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove game %s", g.Name)
	}

	return nil
}

// UpdateEnable is used to update the enable state of a game.
func (g *GameService) UpdateEnable(enabled bool) (bool, error) {
	filter := map[string]string{"name": g.Name}
	currGame, err := g.Get()

	if currGame.Name == "" || err != nil {
		return false, errors.NotFoundf("the game %s", g.Name)
	}

	game := &models.Game{Enabled: &enabled}
	updated, err := game.Update(g.DB, filter)
	if err != nil {
		return false, errors.Errorf("failed to update game %s", err)
	}
	return updated, err
}

func (g *GameService) doesItExist() bool {
	game, err := g.Get()
	if err != nil {
		return false
	}

	return game.Name != ""
}
