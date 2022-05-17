package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func NewRestartCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeployOptions{}
	opts.ApiClient = f.ApiClient
	opts.IO = f.IOStreams
	opts.Restart = true

	restartCmd := &cobra.Command{
		Use:   "restart [project]",
		Short: "Restart a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]

			return runDeploy(opts)
		},
	}

	return restartCmd
}
