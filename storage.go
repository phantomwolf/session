package session

// Provides API for saving/loading/deleting map[string]string objects
type Storage interface {
	Load(key string) (map[string]string, error)
	Save(key string, uid string, data map[string]interface{}) error
	Delete(key string) error
	Exists(key string) bool
}
