package cmd

import (
	"github.com/spf13/cobra"
)

func createRestartCmd() *cobra.Command {
	restartCmd := &cobra.Command{
		Use:   "restart [project]",
		Short: "Restart a project",
		Args:  cobra.ExactArgs(1),
		RunE: withCmdConfig(func(c *CmdConfig) error {
			return checkLoginAndRun(c, Restart, struct{}{})
		}),
	}

	return restartCmd
}

func Restart(c *CmdConfig, _ struct{}) error {
	return Deploy(c, &deployOptions{
		restart: true,
	})
}

func init() {
	rootCmd.AddCommand(createRestartCmd())
}
