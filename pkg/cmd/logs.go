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
	stage        string
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
	logsCmd.Flags().StringVar(&opts.stage, "stage", "runtime", "TODO")

	return logsCmd
}

func Logs(c *CmdConfig, opts *logsOptions) error {
	var deployment *api.Deployment
	var err error

	// Get deployment
	if opts.deploymentID == "" {
		deployment, err = c.client.GetProductionDeployment(c.ctx, c.args[0])
	} else {
		deployment, err = c.client.GetDeployment(c.ctx, uuid.MustParse(opts.deploymentID))
	}
	if err != nil {
		return err
	}

	validStage, logsGetter := logStageToGetter(c, opts.stage, deployment.ID)
	if !validStage {
		return fmt.Errorf("invalid stage name. it must be either \"build\", \"deployment\" or \"runtime\" ")
	}

	if opts.live {
		getStatus := func() (api.DeploymentStatus, error) {
			deployment, err := c.client.GetProductionDeployment(c.ctx, c.args[0])
			if err != nil {
				return deployment.Status, err
			}
			return deployment.Status, nil
		}
		if err := pollLogs(logsGetter, getStatus); err != nil {
			return err
		}
	} else {
		logs, err := logsGetter()
		if err != nil {
			return err
		}
		for _, log := range logs {
			fmt.Println(log.Text)
		}
	}

	return nil
}

func logStageToGetter(c *CmdConfig, stage string, deploymentID uuid.UUID) (valid bool, getter func() ([]api.LogEntry, error)) {
	if stage == "runtime" {
		return true, func() ([]api.LogEntry, error) {
			return c.client.GetRuntimeLogs(c.ctx, deploymentID)
		}
	}

	if stage == "build" {
		return true, func() ([]api.LogEntry, error) {
			return c.client.GetBuildLogs(c.ctx, deploymentID)
		}
	}

	if stage == "deployment" {
		return true, func() ([]api.LogEntry, error) {
			return c.client.GetDeploymentLogs(c.ctx, deploymentID)
		}
	}

	return false, func() ([]api.LogEntry, error) {
		return nil, nil
	}
}

func init() {
	rootCmd.AddCommand(createLogsCmd())
}
