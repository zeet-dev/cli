package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func createLogsCmd() *cobra.Command {
	logsCmd := &cobra.Command{
		Use:   "logs [project]",
		Short: "Logs the output for a given project",
		Args:  cobra.ExactArgs(1),
		RunE: withCmdConfig(func(c *CmdConfig) error {
			return checkLoginAndRun(c, Logs, struct{}{})
		}),
	}

	return logsCmd
}

func Logs(c *CmdConfig, _ struct{}) error {
	project, err := c.client.GetProjectByPath(c.ctx, c.args[0])
	if err != nil {
		return err
	}

	logs, err := c.client.GetProjectLogs(c.ctx, project.ID)
	if err != nil {
		return err
	}

	for _, log := range logs {
		fmt.Print(log)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createLogsCmd())
}
