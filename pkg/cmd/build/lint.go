package build

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/build"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func InitCmds(f *cmdutil.Factory, root *cobra.Command) {
	buildCmd := &cobra.Command{
		Use:   "build <PATH>",
		Short: "Build",
	}

	var target string

	lintCmd := &cobra.Command{
		Use:   "build:lint [OPTIONS] <PATH>",
		Short: "Lint code",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir := args[0]

			return build.LintFiles(workDir, target)
		},
	}

	lintCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target runtime environment")

	root.AddCommand(buildCmd)
	root.AddCommand(lintCmd)
}
