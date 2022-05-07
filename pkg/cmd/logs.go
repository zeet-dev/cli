package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Logs",
	Args:  cobra.ExactArgs(1),
	RunE: withCmdConfig(func(c *CmdConfig) error {
		return checkLoginAndRun(c, Logs)
	}),
}

func Logs(c *CmdConfig) error {
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
	rootCmd.AddCommand(logsCmd)
}
