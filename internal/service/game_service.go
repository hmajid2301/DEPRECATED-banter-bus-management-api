package service

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
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

	t := true
	game := models.Game{
		Name:     g.Name,
		RulesURL: rulesURL,
		Enabled:  &t,
	}

	inserted, err := game.Add(g.DB)
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

	emptyFilter := map[string]string{}
	err := games.Get(g.DB, emptyFilter)
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

// Remove is used to remove a game and all of its questions.
func (g *GameService) Remove() error {
	exists := g.doesItExist()
	if !exists {
		return errors.NotFoundf("the game %s", g.Name)
	}

	questions := models.Questions{}
	filter := map[string]string{"game_name": g.Name}
	_, err := questions.Delete(g.DB, filter)
	if err != nil {
		return err
	}

	filter = map[string]string{"name": g.Name}
	deleted, err := g.DB.Delete("game", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove game %s", g.Name)
	}

	return nil
}

// UpdateEnable is used to update the enable state of a game.
func (g *GameService) UpdateEnable(enabled bool) (bool, error) {
	game, err := g.Get()

	if game.Name == "" || err != nil {
		return false, errors.NotFoundf("The game %s", g.Name)
	}

	newGame := &models.Game{Name: g.Name, RulesURL: game.RulesURL, Enabled: &enabled}
	filter := map[string]string{"name": g.Name}
	updated, err := newGame.Update(g.DB, filter)
	if err != nil {
		return false, errors.Errorf("Failed to update g %s", err)
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
