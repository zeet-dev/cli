package cmdutil

import (
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

type Factory struct {
	IOStreams *iostreams.IOStreams

	ApiClient func() (*api.Client, error)
	Config    func() (config.Config, error)
}
