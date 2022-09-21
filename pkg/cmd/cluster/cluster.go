package cluster

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func NewClusterCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster [command]",
		Short: "Manage clusters",
		Args:  cobra.ExactArgs(1),
	}

	cmd.AddCommand(NewKubeconfigSetCmd(f))
	cmd.AddCommand(NewKubeconfigGetCmd(f))

	return cmd
}
