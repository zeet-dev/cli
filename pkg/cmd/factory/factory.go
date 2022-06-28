package factory

import (
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

func New(version string) *cmdutil.Factory {
	f := &cmdutil.Factory{
		Config:    configFunc(), // No factory dependencies
		IOStreams: ioStreams(),  // No factory dependencies
	}

	f.ApiClient = apiClientFunc(f, version) // Depends on Config, IOStreams, and appVersion

	return f
}

func configFunc() func() (config.Config, error) {
	return func() (config.Config, error) {
		return config.New(), nil
	}
}

func ioStreams() *iostreams.IOStreams {
	return iostreams.System()
}

func apiClientFunc(f *cmdutil.Factory, version string) func() (*api.Client, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil
	}

	return func() (*api.Client, error) {
		return api.New(
			cfg.GetString("api-url"),
			cfg.GetString("auth.access_token"),
			version,
			cfg.GetBool("debug"),
		), nil
	}
}
