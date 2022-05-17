package config

import (
	"github.com/spf13/viper"
)

// Ensure type implements interface
var _ Config = &Live{}

// Live provides a viper-backed config for commands
type Live struct{}

func (c *Live) GetInt(key string) int {
	return viper.GetInt(key)
}

func (c *Live) GetString(key string) string {
	return viper.GetString(key)
}

func (c *Live) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c *Live) Set(key string, val interface{}) error {
	viper.Set(key, val)
	return nil
}

func (c *Live) WriteConfig() error {
	return viper.WriteConfig()
}
