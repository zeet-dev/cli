package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/utils"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy",
	Args:  cobra.ExactArgs(1),
	RunE: withCmdConfig(func(c *CmdConfig) error {
		return checkLoginAndRun(c, Deploy)
	}),
}

func Deploy(c *CmdConfig) error {
	fmt.Printf("Building %s\n", c.args[0])

	project, err := c.client.GetProjectByPath(c.ctx, c.args[0])
	if err != nil {
		return err
	}

	// Build project
	deployment, err := c.client.BuildProject(c.ctx, project.ID, c.cfg.GetString("branch"), c.cfg.GetBool("use-cache"))
	if err != nil {
		return err
	}

	if err := printBuildLogs(c, deployment); err != nil {
		return err
	}

	deployment, err = c.client.GetDeployment(c.ctx, deployment.ID)
	if api.IsBuildSuccess(deployment.Status) {
		fmt.Println(color.GreenString("⛏ ️Build succeeded. Starting deployment..."))
	} else {
		fmt.Println(color.RedString("Build failed"))
		return nil
	}

	if err := printDeploymentLogs(c, deployment); err != nil {
		return err
	}

	deployment, err = c.client.GetDeployment(c.ctx, deployment.ID)
	if api.IsDeploySuccess(deployment.Status) {
		fmt.Printf(color.GreenString("\n🚀 Deployed %s"), c.args[0])

		fmt.Printf(color.GreenString("\n\nPublic Endpoints: %s"), utils.DisplayArray(deployment.Endpoints))

		if deployment.PrivateEndpoint != "" {
			fmt.Printf(color.GreenString("\nPrivate Endpoint: %s"), deployment.PrivateEndpoint)
		}
	} else {
		fmt.Println(color.RedString("Deployment failed"))
	}

	return nil
}

func printBuildLogs(c *CmdConfig, deployment *api.Deployment) error {
	getLogs := func() ([]api.LogEntry, error) {
		return c.client.GetBuildLogs(c.ctx, deployment.ID)
	}
	checkStatus := func() (bool, error) {
		deployment, err := c.client.GetDeployment(c.ctx, deployment.ID)
		if err != nil {
			return false, nil
		}
		return api.IsBuildInProgress(deployment.Status), nil
	}
	if err := printLogs(getLogs, checkStatus); err != nil {
		return err
	}

	return nil
}

func printDeploymentLogs(c *CmdConfig, deployment *api.Deployment) error {
	getLogs := func() ([]api.LogEntry, error) {
		return c.client.GetDeploymentLogs(c.ctx, deployment.ID)
	}
	checkStatus := func() (bool, error) {
		deployment, err := c.client.GetDeployment(c.ctx, deployment.ID)
		if err != nil {
			return false, nil
		}
		return api.IsDeployInProgress(deployment.Status), nil
	}
	if err := printLogs(getLogs, checkStatus); err != nil {
		return err
	}

	return nil
}

func printLogs(getter func() ([]api.LogEntry, error), checker func() (bool, error)) (err error) {
	lastLog := 0

	shouldContinue, err := checker()
	if err != nil {
		return err
	}

	for shouldContinue {
		logs, err := getter()
		if err != nil {
			return err
		}

		for _, log := range logs[lastLog:] {
			fmt.Println(log.Text)
		}
		lastLog = len(logs)

		time.Sleep(time.Second)

		shouldContinue, err = checker()
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	deployCmd.Flags().Bool("use-cache", true, "Enable build cache")
	deployCmd.Flags().StringP("branch", "b", "", "Deploy specific branch")

	viper.BindPFlag("use-cache", deployCmd.Flags().Lookup("use-cache"))
	viper.BindPFlag("branch", deployCmd.Flags().Lookup("branch"))

	rootCmd.AddCommand(deployCmd)
}
