package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

type StatusOpts struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	Project string
}

func NewStatusCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &StatusOpts{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	statusCmd := &cobra.Command{
		Use:   "status [project]",
		Short: "Gets the status for a given project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]

			return runStatus(opts)
		},
	}

	return statusCmd
}

func runStatus(opts *StatusOpts) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	deployment, err := client.GetProductionDeployment(context.Background(), opts.Project)
	if err != nil {
		return err
	}
	status, err := client.GetDeploymentReplicaStatus(context.Background(), deployment.ID)
	if err != nil {
		return err
	}

	var statusMessage string
	switch status.State {
	case "deploy failed":
		statusMessage = color.RedString("DEPLOY FAILED")
	case "build failed":
		statusMessage = color.RedString("BUILD FAILED")
	case "deployed":
		statusMessage = color.GreenString("DEPLOYED")
	default:
		statusMessage = strings.ToUpper(status.State)
	}

	fmt.Fprintf(opts.IO.Out, "Status: %s\n", statusMessage)
	fmt.Fprintf(opts.IO.Out, "Healthy Replicas: [%d/%d]\n", status.ReadyReplicas, status.Replicas)
	return nil
}
