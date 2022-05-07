package config

import "github.com/zeet-dev/cli/pkg/api"

type Provider interface {
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	Set(key string, val interface{}) error

	GetAPIClient(host, accessToken string) *api.Client
	// Subscribe starts the GraphQL Subscriptions listener
	//Subscribe() error
}
