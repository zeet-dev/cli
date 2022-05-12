package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/utils"
)

type deployOptions struct {
	branch   string
	useCache bool
	restart  bool
}

func createDeployCmd() *cobra.Command {
	var opts = &deployOptions{}

	deployCmd := &cobra.Command{
		Use:   "deploy [project]",
		Short: "Deploy a project",
		Args:  cobra.ExactArgs(1),
		RunE: withCmdConfig(func(c *CmdConfig) error {
			return checkLoginAndRun(c, Deploy, opts)
		}),
	}

	deployCmd.Flags().BoolVar(&opts.useCache, "use-cache", true, "Enable build cache")
	deployCmd.Flags().StringVarP(&opts.branch, "branch", "b", "", "Deploy specific branch (defaults to your configured production branch) ")
	deployCmd.Flags().BoolVarP(&opts.restart, "restart", "r", false, "Rerun the latest deployment (this will override use-cache)")

	return deployCmd
}

func Deploy(c *CmdConfig, opts *deployOptions) error {
	project, err := c.client.GetProjectByPath(c.ctx, c.args[0])
	if err != nil {
		return err
	}

	// Build project
	var deployment *api.Deployment

	if c.cfg.GetBool("restart") {
		// Get the branch to restart
		branch := opts.branch
		if branch == "" {
			branch, err = c.client.GetProductionBranch(c.ctx, project.ID)
			if err != nil {
				return err
			}
		}

		deployment, err = c.client.DeployProjectBranch(c.ctx, project.ID, branch, opts.useCache)
		if err != nil {
			return err
		}
	} else {
		deployment, err = c.client.BuildProject(c.ctx, project.ID, opts.branch, opts.useCache)
		if err != nil {
			return err
		}
	}

	deploymentFinished := false
	for !deploymentFinished {
		deployment, err = c.client.GetDeployment(c.ctx, deployment.ID)
		if err != nil {
			return err
		}

		switch deployment.Status {
		// Build
		case api.DeploymentStatusBuildInProgress:
			fmt.Printf("⛏ Building %s...\n", c.args[0])
			if err := printBuildLogs(c, deployment); err != nil {
				return err
			}
			break
		case api.DeploymentStatusBuildSucceeded:
			fmt.Println(color.GreenString("⛏ Build complete\n"))
			break
		case api.DeploymentStatusBuildFailed:
			fmt.Println(color.RedString("Build failed\n"))
			deploymentFinished = true
			break
		case api.DeploymentStatusBuildAborted:
			fmt.Println(color.RedString("Build aborted\n"))
			deploymentFinished = true
			break
		case api.DeploymentStatusDeployStopped:
			fmt.Println(color.RedString("Build stopped\n"))
			break

		// Deployment
		case api.DeploymentStatusDeployInProgress:
			fmt.Printf("Deploying %s...\n", c.args[0])
			if err := printDeploymentLogs(c, deployment); err != nil {
				return err
			}
			break
		case api.DeploymentStatusDeploySucceeded:
			printDeploymentSummary(c, deployment)
			deploymentFinished = true
			break
		case api.DeploymentStatusDeployFailed:
			fmt.Println(color.RedString("Deploy failed\n"))
			deploymentFinished = true
			break
		}
	}

	return nil
}

func printBuildLogs(c *CmdConfig, deployment *api.Deployment) error {
	getLogs := func() ([]api.LogEntry, error) {
		return c.client.GetBuildLogs(c.ctx, deployment.ID)
	}
	getStatus := func() (api.DeploymentStatus, error) {
		deployment, err := c.client.GetDeployment(c.ctx, deployment.ID)
		if err != nil {
			return deployment.Status, err
		}
		return deployment.Status, nil
	}
	if err := printLogs(getLogs, getStatus); err != nil {
		return err
	}

	return nil
}

func printDeploymentLogs(c *CmdConfig, deployment *api.Deployment) error {
	getLogs := func() ([]api.LogEntry, error) {
		return c.client.GetDeploymentLogs(c.ctx, deployment.ID)
	}
	getStatus := func() (api.DeploymentStatus, error) {
		deployment, err := c.client.GetDeployment(c.ctx, deployment.ID)
		if err != nil {
			return deployment.Status, err
		}
		return deployment.Status, nil
	}
	if err := printLogs(getLogs, getStatus); err != nil {
		return err
	}

	return nil
}

func printDeploymentSummary(c *CmdConfig, deployment *api.Deployment) {
	fmt.Printf(color.GreenString("\n🚀 Deployed %s"), c.args[0])
	fmt.Printf(color.GreenString("\n\nPublic Endpoints: \n%s"), utils.DisplayArray(deployment.Endpoints))
	if deployment.PrivateEndpoint != "" {
		fmt.Printf(color.GreenString("\nPrivate Endpoint: %s"), deployment.PrivateEndpoint)
	}
}

func init() {
	rootCmd.AddCommand(createDeployCmd())
}
