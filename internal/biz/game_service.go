package biz

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/games"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/biz/models"
)

// GameService is struct data required by all g service functions.
type GameService struct {
	DB models.Repository
}

// Add is used to add a new game.
func (g *GameService) Add(name string, rulesURL string) error {
	exists := g.doesItExist(name)
	if exists {
		return errors.AlreadyExistsf("The game %s", name)
	}

	game, err := games.GetGame(name)
	if err != nil {
		return err
	}

	newGame := game.NewGame(rulesURL)
	inserted, err := newGame.Add(g.DB)
	if !inserted {
		return errors.Errorf("Failed to add the new game %s", name)
	}
	return err
}

// Get is used to get specific information about a game.
func (g *GameService) Get(name string) (*models.Game, error) {
	var (
		filter = map[string]string{"name": name}
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
	for _, g := range games {
		if enabled != nil && *enabled != *g.Enabled {
			continue
		}
		gameNames = append(gameNames, g.Name)
	}
	return gameNames, nil
}

// Remove is used to remove a game.
func (g *GameService) Remove(name string) error {
	exists := g.doesItExist(name)
	if !exists {
		return errors.NotFoundf("the game %s", name)
	}

	filter := map[string]string{"name": name}
	deleted, err := g.DB.Delete("game", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove game %s", name)
	}

	return nil
}

// UpdateEnable is used to update the enable state of a game.
func (g *GameService) UpdateEnable(name string, enabled bool) (bool, error) {
	filter := map[string]string{"name": name}
	game, err := g.Get(name)

	if game.Name == "" || err != nil {
		return false, errors.NotFoundf("The g %s", name)
	}

	newGame := &models.Game{Name: name, RulesURL: game.RulesURL, Enabled: &enabled}
	updated, err := newGame.Update(g.DB, filter)
	if err != nil {
		return false, errors.Errorf("Failed to update g %s", err)
	}
	return updated, err
}

func (g *GameService) doesItExist(name string) bool {
	game, err := g.Get(name)
	if err != nil {
		return false
	}

	return game.Name != ""
}
