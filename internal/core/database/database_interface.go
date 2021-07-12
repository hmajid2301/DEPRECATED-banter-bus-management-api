package database

type Documents interface {
	Add(db Database) error
	Get(db Database, filter map[string]interface{}) error
	GetWithLimit(db Database, filter map[string]interface{}, limit int64) error
	ToInterface() []interface{}
	Delete(db Database, filter map[string]interface{}) (bool, error)
}

type Document interface {
	Add(db Database) (bool, error)
	Get(db Database, filter map[string]interface{}) error
	Update(db Database, filter map[string]interface{}) (bool, error)
}

type NewSubDocument interface {
	AddToList(db Database, filter map[string]interface{}) (bool, error)
}

type UpdateSubDocument interface {
	Add(db Database, filter map[string]interface{}) (bool, error)
	Remove(db Database, filter map[string]interface{}) (bool, error)
}

type SubDocument interface {
	RemoveFromList(db Database, filter map[string]interface{}) (bool, error)
}

type Database interface {
	Ping() bool
	Insert(collectionName string, document Document) (bool, error)
	InsertMultiple(collectionName string, documents Documents) error
	Get(collectionName string, filter map[string]interface{}, document Document) error
	GetAll(collectionName string, filter map[string]interface{}, documents Documents) error
	GetWithLimit(collectionName string, filter map[string]interface{}, limit int64, documents Documents) error
	GetRandom(collectionName string, filter map[string]interface{}, limit int64, documents Documents) error
	GetUniqueValues(
		collectionName string,
		filter map[string]interface{},
		fieldName string,
	) ([]string, error)
	GetUniqueKeys(collectionName string, filter map[string]interface{}, fieldName string) ([]string, error)
	Delete(collectionName string, filter map[string]interface{}) (bool, error)
	DeleteAll(collectionName string, filter map[string]interface{}) (bool, error)
	RemoveCollection(collectionName string) error
	Update(collectionName string, filter map[string]interface{}, document Document) (bool, error)
	UpdateObject(collectionName string, filter map[string]interface{}, subDocument UpdateSubDocument) (bool, error)
	RemoveObject(collectionName string, filter map[string]interface{}, subDocument UpdateSubDocument) (bool, error)
	AppendToList(collectionName string, filter map[string]interface{}, subDocument NewSubDocument) (bool, error)
	RemoveFromList(collectionName string, filter map[string]interface{}, subDocument SubDocument) (bool, error)
}
