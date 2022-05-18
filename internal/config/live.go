package config

import (
	"github.com/spf13/viper"
)

// Ensure type implements interface
var _ Config = &Live{}

// Live provides a viper-backed config for commands
type Live struct {
	v *viper.Viper
}

func (c *Live) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *Live) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}

func (c *Live) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Live) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *Live) Set(key string, val interface{}) error {
	c.v.Set(key, val)
	return nil
}

func (c *Live) WriteConfig() error {
	return c.v.WriteConfig()
}
