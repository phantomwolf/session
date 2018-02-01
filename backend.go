package session

type Backend interface {
	Load(key string) (map[string]string, error)
	Save(key string, data map[string]string) error
	Delete(key string) error
}
