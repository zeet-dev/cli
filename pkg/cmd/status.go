package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func createStatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status [project]",
		Short: "Gets the status for a given project",
		Args:  cobra.ExactArgs(1),
		RunE: withCmdConfig(func(c *CmdConfig) error {
			return checkLoginAndRun(c, Status, struct{}{})
		}),
	}

	return statusCmd
}

func Status(c *CmdConfig, _ struct{}) error {
	deployment, err := c.client.GetProductionDeployment(c.ctx, c.args[0])
	if err != nil {
		return err
	}
	status, err := c.client.GetDeploymentReplicaStatus(c.ctx, deployment.ID)
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

	fmt.Printf("Status: %s\n", statusMessage)
	fmt.Printf("Healthy Replicas: [%d/%d]\n", status.ReadyReplicas, status.Replicas)
	return nil
}

func init() {
	rootCmd.AddCommand(createStatusCmd())
}
