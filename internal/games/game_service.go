package games

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/questions"
)

type GameService struct {
	DB   database.Database
	Name string
}

func (g *GameService) Add(rulesURL string) error {
	exists := g.doesItExist()
	if exists {
		return errors.AlreadyExistsf("The game %s", g.Name)
	}

	t := true
	game := Game{
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

func (g *GameService) Get() (*Game, error) {
	var (
		filter = map[string]string{"name": g.Name}
		game   = &Game{}
	)

	err := game.Get(g.DB, filter)
	if err != nil {
		return &Game{}, err
	}

	return game, nil
}

func (g *GameService) GetAll(enabled *bool) ([]string, error) {
	games := Games{}

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

func (g *GameService) Remove() error {
	exists := g.doesItExist()
	if !exists {
		return errors.NotFoundf("the game %s", g.Name)
	}

	questions := questions.Questions{}
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

func (g *GameService) UpdateEnable(enabled bool) (bool, error) {
	game, err := g.Get()

	if game.Name == "" || err != nil {
		return false, errors.NotFoundf("The game %s", g.Name)
	}

	updateGame := &Game{Name: g.Name, RulesURL: game.RulesURL, Enabled: &enabled}
	filter := map[string]string{"name": g.Name}
	updated, err := updateGame.Update(g.DB, filter)
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
