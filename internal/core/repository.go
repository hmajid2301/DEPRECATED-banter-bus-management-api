package core

// Repository interface defines all the methods required by a repository
type Repository interface {
	Ping() bool
	Insert(collectionName string, objectToInsert interface{}) (bool, error)
	InsertMultiple(collectionName string, objectToInsert []interface{}) error
	Get(collectionName string, filter interface{}, unmarshalModel interface{}) error
	GetAll(collectionName string, unmarshalModel interface{}) error
	Delete(collectionName string, filter interface{}) (bool, error)
	RemoveCollection(collectionName string) error
	UpdateEntry(collectionName string, filter interface{}, objectToUpdate interface{}) (bool, error)
	RemoveEntry(collectionName string, filter interface{}, objectToRemove interface{}) (bool, error)
	AppendToEntry(collectionName string, filter interface{}, objectToAppend interface{}) (bool, error)
	RemoveFromEntry(collectionName string, filter interface{}, objectToRemove interface{}) (bool, error)
}
