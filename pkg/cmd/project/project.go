package project

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func NewProjectCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "project [command]",
		Short: "Manage projects",
		Args: cobra.ExactArgs(1),
	}

	cmd.AddCommand(NewProjectListCmd(f))
	cmd.AddCommand(NewProjectShowCmd(f))

	return cmd
}
