package core

import (
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/dbmodels"

	"github.com/juju/errors"
)

// AddGameType is used add a new game type.
func AddGameType(name string, rulesURL string) error {
	exists := doesGameExist(name)
	if exists {
		return errors.AlreadyExistsf("The game type %s", name)
	}

	t := true
	var newGame = dbmodels.Game{
		Name:     name,
		RulesURL: rulesURL,
		Questions: &dbmodels.Question{
			Rounds: &dbmodels.Rounds{
				One:   []string{},
				Two:   []string{},
				Three: []string{},
			},
		},
		Enabled: &t,
	}

	inserted, err := database.Insert("game", newGame)
	if !inserted {
		return errors.Errorf("Failed to add new game type %s", name)
	}
	return err
}

// GetGameType is used to specific information about a game type.
func GetGameType(name string) (*dbmodels.Game, error) {
	var (
		filter = dbmodels.Game{Name: name}
		game   *dbmodels.Game
	)

	err := database.Get("game", filter, &game)
	if err != nil {
		return &dbmodels.Game{}, err
	}

	return game, nil
}

// GetAllGameTypes is used to get all game types.
func GetAllGameTypes(filter *bool) ([]string, error) {
	games := []*dbmodels.Game{}

	err := database.GetAll("game", &games)
	if err != nil {
		return []string{}, err
	}

	var gameNames []string
	for _, game := range games {
		if filter != nil && *filter != *game.Enabled {
			continue
		}
		gameNames = append(gameNames, game.Name)
	}
	return gameNames, nil
}

// RemoveGameTypes is used to remove a game type.
func RemoveGameTypes(name string) error {
	exists := doesGameExist(name)
	if !exists {
		return errors.NotFoundf("The game type %s", name)
	}

	var filter = dbmodels.Game{Name: name}
	deleted, err := database.Delete("game", filter)
	if !deleted || err != nil {
		return errors.Errorf("Failed to remove game type %s", name)
	}

	return nil
}

// UpdateEnableGameType is used to update the enable state of a game type.
func UpdateEnableGameType(name string, enabled bool) (bool, error) {
	var filter = dbmodels.Game{Name: name}
	game, _ := GetGameType(name)

	if game.Name == "" {
		return false, errors.NotFoundf("The game type %s", name)
	} else if *game.Enabled == enabled {
		return false, errors.AlreadyExistsf("The game type is already enabled: %s", enabled)
	}

	var update = dbmodels.Game{Name: name, Enabled: &enabled}
	updated, err := database.UpdateEntry("game", filter, update)
	if err != nil {
		return false, errors.Errorf("Failed to update game %s", err)
	}
	return updated, err
}

func doesGameExist(name string) bool {
	game, _ := GetGameType(name)
	return game.Name != ""
}
