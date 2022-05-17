package config

type Config interface {
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	Set(key string, val interface{}) error

	WriteConfig() error
}

// New returns a new Live config
func New() Config {
	return &Live{}
}
