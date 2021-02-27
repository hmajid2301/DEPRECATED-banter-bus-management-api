package database

// Documents are a list of document(s) in the database
type Documents interface {
	Add(db Database) error
	Get(db Database) error
	ToInterface() []interface{}
}

// Document is a single "item" in the database i.e. a single user.
type Document interface {
	Add(db Database) (bool, error)
	Get(db Database, filter map[string]string) error
	Update(db Database, filter map[string]string) (bool, error)
}

// NewSubDocument are part of document(s), being added in.
type NewSubDocument interface {
	AddToList(db Database, filter map[string]string) (bool, error)
}

// SubDocument are part of document(s), usually as sub-objects.
type SubDocument interface {
	RemoveFromList(db Database, filter map[string]string) (bool, error)
}

// SubDocuments are a list of sub-objects of a document.
type SubDocuments interface {
	Get(db Database, filter map[string]string, parentField string, condition []string) error
}

// Database defines the contract when interacting with a repository (i.e. MongoDB).
type Database interface {
	Ping() bool
	Insert(collectionName string, document Document) (bool, error)
	InsertMultiple(collectionName string, documents Documents) error
	Get(collectionName string, filter map[string]string, document Document) error
	GetAll(collectionName string, documents Documents) error
	Delete(collectionName string, filter map[string]string) (bool, error)
	RemoveCollection(collectionName string) error
	Update(collectionName string, filter map[string]string, document Document) (bool, error)
	GetSubObject(
		collectionName string,
		filter map[string]string,
		parentField string,
		condition []string,
		subDocument SubDocuments,
	) error
	UpdateObject(collectionName string, filter map[string]string, update map[string]interface{}) (bool, error)
	RemoveObject(collectionName string, filter map[string]string, remove map[string]interface{}) (bool, error)
	AppendToList(collectionName string, filter map[string]string, subDocument NewSubDocument) (bool, error)
	RemoveFromList(collectionName string, filter map[string]string, subDocument SubDocument) (bool, error)
}
