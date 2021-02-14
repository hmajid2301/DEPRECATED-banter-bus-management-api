package models

// Documents are a list of document(s) in the database.
type Documents interface {
	Add(db Repository) error
	Get(db Repository) error
	ToInterface() []interface{}
}

// Document is a single "item" in the database i.e. a single user.
type Document interface {
	Add(db Repository) (bool, error)
	Get(db Repository, filter map[string]string) error
	Update(db Repository, filter map[string]string) (bool, error)
}

// NewSubDocument are part of document(s), being added in.
type NewSubDocument interface {
	AddToList(db Repository, filter map[string]string) (bool, error)
}

// SubDocument are part of document(s), usually as sub-objects.
type SubDocument interface {
	RemoveFromList(db Repository, filter map[string]string) (bool, error)
}

// SubDocuments are a list of sub-objects of a document.
type SubDocuments interface {
	Get(db Repository, filter map[string]string, parentField string, condition []string) error
}

// Repository defines the contract when interacting with a repository (i.e. MongoDB).
type Repository interface {
	Ping() bool
	Insert(collectionName string, objectToInsert Document) (bool, error)
	InsertMultiple(collectionName string, objectToInsert Documents) error
	Get(collectionName string, filter map[string]string, objectToGet Document) error
	GetAll(collectionName string, objectToGet Documents) error
	Delete(collectionName string, filter map[string]string) (bool, error)
	RemoveCollection(collectionName string) error
	Update(collectionName string, filter map[string]string, objectToUpdate Document) (bool, error)
	GetSubObject(
		collectionName string,
		filter map[string]string,
		parentField string,
		condition []string,
		model SubDocuments,
	) error
	UpdateObject(collectionName string, filter map[string]string, objectToUpdate map[string]interface{}) (bool, error)
	RemoveObject(collectionName string, filter map[string]string, objectToRemove map[string]interface{}) (bool, error)
	AppendToList(collectionName string, filter map[string]string, objectToAppend NewSubDocument) (bool, error)
	RemoveFromList(collectionName string, filter map[string]string, objectToRemove SubDocument) (bool, error)
}
