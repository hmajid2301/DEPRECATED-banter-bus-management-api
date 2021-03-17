package database

// Documents are a list of document(s) in the database
type Documents interface {
	Add(db Database) error
	Get(db Database, filter map[string]string) error
	ToInterface() []interface{}
	Delete(db Database, filter map[string]string) (bool, error)
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

// UpdateSubDocument are part of document(s), being added in.
type UpdateSubDocument interface {
	Add(db Database, filter map[string]string) (bool, error)
	Remove(db Database, filter map[string]string) (bool, error)
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
	GetAll(collectionName string, filter map[string]string, documents Documents) error
	GetUnique(
		collectionName string,
		filter map[string]string,
		fieldName string,
	) ([]string, error)
	Delete(collectionName string, filter map[string]string) (bool, error)
	DeleteAll(collectionName string, filter map[string]string) (bool, error)
	RemoveCollection(collectionName string) error
	Update(collectionName string, filter map[string]string, document Document) (bool, error)
	GetSubObject(
		collectionName string,
		filter map[string]string,
		parentField string,
		condition []string,
		subDocument SubDocuments,
	) error
	UpdateObject(collectionName string, filter map[string]string, subDocument UpdateSubDocument) (bool, error)
	RemoveObject(collectionName string, filter map[string]string, subDocument UpdateSubDocument) (bool, error)
	AppendToList(collectionName string, filter map[string]string, subDocument NewSubDocument) (bool, error)
	RemoveFromList(collectionName string, filter map[string]string, subDocument SubDocument) (bool, error)
}
