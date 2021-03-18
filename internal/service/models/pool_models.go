package models

import (
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

// Pool struct holds the data for all their own questions
type Pool struct {
	Username     string `bson:"username"`
	PoolName     string `bson:"pool_name"     json:"pool_name"`
	GameName     string `bson:"game_name"     json:"game_name"`
	LanguageCode string `bson:"language_code" json:"language_code"`
	Privacy      string `bson:"privacy"`
}

// Add is used to add a pool.
func (pool *Pool) Add(db database.Database) (bool, error) {
	inserted, err := db.Insert("pool", pool)
	return inserted, err
}

// Get is used to get a pool.
func (pool *Pool) Get(db database.Database, filter map[string]string) error {
	err := db.Get("pool", filter, pool)
	return err
}

// Update is used to update a pool.
func (pool *Pool) Update(db database.Database, filter map[string]string) (bool, error) {
	updated, err := db.Update("pool", filter, pool)
	return updated, err
}

// Pools is a list of pool(s).
type Pools []Pool

// Add method adds (a list of) pool at once.
func (pools *Pools) Add(db database.Database) error {
	err := db.InsertMultiple("pool", pools)
	return err
}

// Get method gets all the pools.
func (pools *Pools) Get(db database.Database, filter map[string]string) error {
	err := db.GetAll("pool", filter, pools)
	return err
}

// Delete is used to delete a list of pools that match a filter.
func (pools Pools) Delete(db database.Database, filter map[string]string) (bool, error) {
	deleted, err := db.DeleteAll("pool", filter)
	return deleted, err
}

// ToInterface converts pools (list of pools) into a list of interfaces, required by GetAll MongoDB.
func (pools Pools) ToInterface() []interface{} {
	interfaceObject := make([]interface{}, len(pools))
	for i, item := range pools {
		interfaceObject[i] = item
	}
	return interfaceObject
}
