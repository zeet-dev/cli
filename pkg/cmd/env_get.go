package cmd

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
	"github.com/zeet-dev/cli/pkg/utils"
)

type EnvGetOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	Project string
	Keys    []string
}

func NewEnvGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &EnvGetOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use: "env:get [project] [name?]",
		Example: heredoc.Doc(`
			$ zeet env:get zeet-demo/zeet-demo-node-sample DB_USERNAME
			$ zeet env:get zeet-demo/zeet-demo-node-sample DB_USERNAME DB_PASSWORD
			$ zeet env:get zeet-demo/zeet-demo-node-sample
		`),
		Short: "Retrieve an environment variable for a project. Pass only the project name to get a list of all variables.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]
			opts.Keys = args[1:]

			return runEnvGet(opts)
		},
	}

	return cmd
}

func runEnvGet(opts *EnvGetOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	repo, err := client.GetProjectRepo(context.Background(), opts.Project)
	if err != nil {
		return err
	}

	vars, err := client.GetEnvVars(context.Background(), repo.ID)
	if err != nil {
		return err
	}

	if len(opts.Keys) == 0 {
		fmt.Fprintln(opts.IO.Out, utils.DisplayMap(vars))

	} else {
		for k, v := range vars {
			if utils.SliceContains(opts.Keys, k) {
				fmt.Fprintf(opts.IO.Out, "%s=%s\n", k, v)
			}
		}
	}

	return nil
}
