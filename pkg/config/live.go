package config

import (
	"github.com/spf13/viper"
	"github.com/zeet-dev/cli/pkg/api"
)

// Ensure type implements interface
var _ Provider = &Live{}

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

func (c *Live) GetAPIClient(host, accessToken string) *api.Client {
	return api.New(host, accessToken)
}

func (c *Live) Set(key string, val interface{}) error {
	viper.Set(key, val)
	return nil
}
