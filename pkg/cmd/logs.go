package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		projectPath := args[0]
		project, err := api.GetProject(ctx, projectPath)
		if err != nil {
			return err
		}

		logs, err := api.GetProjectLogs(ctx, string(project.ID))
		if err != nil {
			return err
		}

		for _, log := range logs {
			fmt.Println(log)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
