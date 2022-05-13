package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
)

type logsOptions struct {
	live         bool
	deploymentID string
}

func createLogsCmd() *cobra.Command {
	var opts = &logsOptions{}

	logsCmd := &cobra.Command{
		Use:   "logs [project]",
		Short: "Logs the output for a given project",
		Args:  cobra.ExactArgs(1),
		RunE: withCmdConfig(func(c *CmdConfig) error {
			return checkLoginAndRun(c, Logs, opts)
		}),
	}

	logsCmd.Flags().BoolVarP(&opts.live, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVarP(&opts.deploymentID, "deployment", "d", "", "The ID of the deployment to get logs for (not respected for serverless)")

	return logsCmd
}

func Logs(c *CmdConfig, opts *logsOptions) error {
	project, err := c.client.GetProjectByPath(c.ctx, c.args[0])
	if err != nil {
		return err
	}

	var logsGetter func() ([]api.LogEntry, error)
	if opts.deploymentID == "" {
		logsGetter = func() ([]api.LogEntry, error) {
			return c.client.GetProjectLogs(c.ctx, project.ID)
		}
	} else {
		logsGetter = func() ([]api.LogEntry, error) {
			return c.client.GetDeploymentLogs(c.ctx, uuid.MustParse(opts.deploymentID))
		}
	}

	if opts.live && opts.deploymentID == "" {
		getStatus := func() (api.DeploymentStatus, error) {
			deployment, err := c.client.GetProductionDeployment(c.ctx, c.args[0])
			if err != nil {
				return deployment.Status, err
			}
			return deployment.Status, nil
		}
		if err := printLogs(logsGetter, getStatus); err != nil {
			return err
		}
	} else {
		logs, err := c.client.GetDeploymentLogs(c.ctx, uuid.MustParse(opts.deploymentID))
		if err != nil {
			return err
		}
		for _, log := range logs {
			fmt.Println(log.Text)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createLogsCmd())
}
