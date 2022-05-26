package cmd

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
	"github.com/zeet-dev/cli/pkg/utils"
)

type LogsOptions struct {
	IO        *iostreams.IOStreams
	Config    func() (config.Config, error)
	ApiClient func() (*api.Client, error)

	Live         bool
	DeploymentID string
	Stage        string
	Project      string
}

func NewLogsCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &LogsOptions{}
	opts.IO = f.IOStreams
	opts.Config = f.Config
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:   "logs [project]",
		Short: "Logs the output for a given project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]

			return runLogs(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Live, "follow", "f", false, "Follow log output")
	cmd.Flags().StringVarP(&opts.DeploymentID, "deployment", "d", "", "The ID of the deployment to get logs for (not respected for serverless)")
	cmd.Flags().StringVar(&opts.Stage, "stage", "runtime", "The deployment stage to get the logs for (build, deployment, or runtime)")

	return cmd
}

func runLogs(opts *LogsOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	path, err := utils.ToProjectPath(client, opts.Project)
	if err != nil {
		return err
	}

	// Get deployment
	var deployment *api.Deployment
	if opts.DeploymentID == "" {
		deployment, err = client.GetProductionDeployment(context.Background(), path)
	} else {
		deployment, err = client.GetDeployment(context.Background(), uuid.MustParse(opts.DeploymentID))
	}
	if err != nil {
		return err
	}

	validStage, logsGetter := logStageToGetter(client, opts.Stage, deployment.ID)
	if !validStage {
		return fmt.Errorf("invalid stage name. it must be either \"build\", \"deployment\" or \"runtime\" ")
	}

	if opts.Live {
		getStatus := func() (api.DeploymentStatus, error) {
			deployment, err := client.GetProductionDeployment(context.Background(), path)
			if err != nil {
				return deployment.Status, err
			}
			return deployment.Status, nil
		}
		if err := utils.PollLogs(logsGetter, getStatus, opts.IO.Out); err != nil {
			return err
		}
	} else {
		logs, err := logsGetter()
		if err != nil {
			return err
		}
		for _, log := range logs {
			fmt.Fprintln(opts.IO.Out, log.Text)
		}
	}

	return nil
}

func logStageToGetter(client *api.Client, stage string, deploymentID uuid.UUID) (valid bool, getter func() ([]api.LogEntry, error)) {
	if stage == "runtime" {
		return true, func() ([]api.LogEntry, error) {
			return client.GetRuntimeLogs(context.Background(), deploymentID)
		}
	}

	if stage == "build" {
		return true, func() ([]api.LogEntry, error) {
			return client.GetBuildLogs(context.Background(), deploymentID)
		}
	}

	if stage == "deployment" {
		return true, func() ([]api.LogEntry, error) {
			return client.GetDeploymentLogs(context.Background(), deploymentID)
		}
	}

	return false, func() ([]api.LogEntry, error) {
		return nil, nil
	}
}
