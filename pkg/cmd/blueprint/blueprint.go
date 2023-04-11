package blueprint

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func NewBlueprintCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blueprint [command]",
		Short: "Manage blueprints",
		Args:  cobra.ExactArgs(1),
	}

	cmd.AddCommand(NewBlueprintListCmd(f))
	cmd.AddCommand(NewBlueprintShowCmd(f))

	return cmd
}
