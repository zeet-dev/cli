package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
	"github.com/zeet-dev/cli/pkg/utils"
)

type DeployOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	Image    string
	Branch   string
	Project  string
	UseCache bool
	Restart  bool
	Follow   bool
}

func NewDeployCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &DeployOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	deployCmd := &cobra.Command{
		Use:   "deploy [project]",
		Short: "Deploy a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]

			return runDeploy(opts)
		},
	}

	deployCmd.Flags().BoolVar(&opts.UseCache, "use-cache", true, "Enable build cache")
	deployCmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Deploy specific Branch (defaults to your configured production branch) ")
	deployCmd.Flags().StringVarP(&opts.Image, "image", "i", "", "The Docker image to use for this deployment (if used with --branch, only the branch's image will be updated)")
	deployCmd.Flags().BoolVarP(&opts.Follow, "follow", "f", true, "Follow the deployment logs. If false, the deployment will be started then the command will exit")

	return deployCmd
}

func runDeploy(opts *DeployOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	path, err := client.ToProjectPath(opts.Project)
	if err != nil {
		return err
	}
	project, err := client.GetProjectByPath(context.Background(), path)
	if err != nil {
		return err
	}

	// Build project
	var deployment *api.Deployment

	if opts.Restart {
		// Get the branch to Restart
		branch := opts.Branch
		if branch == "" {
			branch, err = client.GetProductionBranch(context.Background(), project.ID)
			if err != nil {
				return err
			}
		}

		deployment, err = client.DeployProjectBranch(context.Background(), project.ID, branch, opts.UseCache)
		if err != nil {
			return err
		}
	} else if opts.Image != "" {
		deployment, err = updateImage(client, project, path, opts.Image, opts.Branch)
		if err != nil {
			return err
		}
	} else {
		deployment, err = client.BuildProject(context.Background(), project.ID, opts.Branch, opts.UseCache)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(opts.IO.Out, "Deploy started...")
	dashboard := fmt.Sprintf("https://zeet.co/repo/%s/deployments/%s", project.ID, deployment.ID)
	fmt.Fprintln(opts.IO.Out, "Dashboard: "+dashboard)
	if !opts.Follow {
		return nil
	}

	for {
		deployment, err = client.GetDeployment(context.Background(), deployment.ID)
		if err != nil {
			return err
		}
		fullStatus, err := client.GetDeploymentStatus(context.Background(), deployment.ID)
		if err != nil {
			return err
		}
		if fullStatus.State == "error" {
			return fmt.Errorf("Deploy failed. %s\n", fullStatus.ErrorMessage)
		}

		switch deployment.Status {
		// Build
		case api.DeploymentStatusBuildInProgress:
			fmt.Fprintf(opts.IO.Out, "⛏ Building %s...\n", path)
			if err := printBuildLogs(client, deployment, opts.IO.Out); err != nil {
				return err
			}
			break
		case api.DeploymentStatusBuildSucceeded:
			fmt.Fprintf(opts.IO.Out, color.GreenString("⛏ Build complete\n"))
			return nil
		case api.DeploymentStatusBuildFailed:
			return fmt.Errorf("Build failed\n")
		case api.DeploymentStatusBuildAborted:
			return fmt.Errorf("Build aborted\n")
		case api.DeploymentStatusDeployStopped:
			return fmt.Errorf("Build stopped\n")

		// Deployment
		case api.DeploymentStatusDeployInProgress:
			fmt.Fprintf(opts.IO.Out, "Deploying %s...\n", path)
			if err := printDeploymentLogs(client, deployment, opts.IO.Out); err != nil {
				return err
			}
			break
		case api.DeploymentStatusDeploySucceeded:
			printDeploymentSummary(deployment, path, opts.IO.Out)
			return nil
		case api.DeploymentStatusDeployFailed:
			return fmt.Errorf("Deploy failed. %s\n", deployment.ErrorMessage)
		}
	}
}

func printBuildLogs(client *api.Client, deployment *api.Deployment, out io.Writer) error {
	getLogs := func() ([]api.LogEntry, error) {
		return client.GetBuildLogs(context.Background(), deployment.ID)
	}
	shouldContinue := func() (bool, error) {
		deployment, err := client.GetDeployment(context.Background(), deployment.ID)
		if err != nil {
			return false, err
		}
		return api.IsBuildInProgress(deployment.Status), nil
	}
	if err := pollLogs(getLogs, shouldContinue, out); err != nil {
		return err
	}

	return nil
}

func printDeploymentLogs(client *api.Client, deployment *api.Deployment, out io.Writer) error {
	getLogs := func() ([]api.LogEntry, error) {
		return client.GetDeploymentLogs(context.Background(), deployment.ID)
	}
	getStatus := func() (bool, error) {
		deployment, err := client.GetDeployment(context.Background(), deployment.ID)
		if err != nil {
			return false, err
		}
		return api.IsDeployInProgress(deployment.Status), nil
	}
	if err := pollLogs(getLogs, getStatus, out); err != nil {
		return err
	}

	return nil
}

func printDeploymentSummary(deployment *api.Deployment, project string, out io.Writer) {
	fmt.Fprintf(out, color.GreenString("\n🚀 Deployed %s"), project)
	fmt.Fprintf(out, color.GreenString("\n\nPublic Endpoints: \n%s"), utils.DisplayArray(deployment.Endpoints))
	if deployment.PrivateEndpoint != "" {
		fmt.Printf(color.GreenString("\nPrivate Endpoint: %s\n"), deployment.PrivateEndpoint)
	}
}

// updateImage updates a project's Docker image. If branch is not empty, it will update a branch's image instead
func updateImage(client *api.Client, project *api.Project, projectPath string, image, branch string) (deployment *api.Deployment, err error) {
	// Update
	if branch == "" {
		err = client.UpdateProject(context.Background(), project.ID, image)
		if err != nil {
			return
		}
	} else {
		err = client.UpdateBranch(context.Background(), project.ID, image, branch, true)
		if err != nil {
			return
		}
	}

	// Get the resulting deployment
	if branch == "" {
		deployment, err = client.GetProductionDeployment(context.Background(), projectPath)
	} else {
		deployment, err = client.GetLatestDeployment(context.Background(), projectPath, branch)
	}
	return
}
