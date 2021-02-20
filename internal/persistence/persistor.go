package persistence

type Persistor interface {
	Store(key string, data interface{}) error
	Load(key string) (interface{}, error)
}

type FSPersistor struct {
	
}