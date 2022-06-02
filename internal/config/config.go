package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config interface {
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	Set(key string, val interface{}) error

	WriteConfig() error
}

// New returns a new Live config
func New() Config {
	return &Live{viper.GetViper()}
}

// TODO make NewConfig like this?
func NewState() (Config, error) {
	v := viper.New()
	dir, err := stateDir()
	if err != nil {
		return nil, err
	}

	ch := filepath.Join(dir, "state.yaml")
	v.SetConfigFile(ch)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
			// No problem, the config file will be created later
		} else {
			return nil, err
		}
	}

	return &Live{v}, nil
}

func ConfigDir() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	ch := filepath.Join(cfgDir, "zeet")
	err = os.MkdirAll(ch, os.ModePerm)
	if err := os.MkdirAll(ch, os.ModePerm); err != nil {
		return "", err
	}

	return ch, nil
}

// State path precedence
// 1. XDG_STATE_HOME
// 2. LocalAppData (windows only)
// 3. HOME
// (from gh cli)
func stateDir() (string, error) {
	var path string
	if a := os.Getenv("XDG_STATE_HOME"); a != "" {
		path = filepath.Join(a, "zeet")
	} else if b := os.Getenv("LOCAL_APP_DATA"); runtime.GOOS == "windows" && b != "" {
		path = filepath.Join(b, "Zeet CLI")
	} else {
		c, _ := os.UserHomeDir()
		path = filepath.Join(c, ".local", "state", "zeet")
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	return path, nil
}
