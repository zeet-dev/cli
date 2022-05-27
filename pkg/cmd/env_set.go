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
	"github.com/zeet-dev/cli/pkg/utils"
)

type EnvSetOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	Project string
	Vars    []string
}

func NewEnvSetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &EnvSetOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:     "env:set [project] [name=value]",
		Example: "$ zeet env:set zeet-demo/zeet-demo-node-sample DB_USERNAME=test DB_PASSWORD=ilovezeet123!",
		Short:   "Add or modify an environment variable for a project",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]
			opts.Vars = args[1:]

			return runEnvSet(opts)
		},
	}

	return cmd
}

func runEnvSet(opts *EnvSetOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	path, err := utils.ToProjectPath(client, opts.Project)
	if err != nil {
		return err
	}

	repo, err := client.GetProjectRepo(context.Background(), path)
	if err != nil {
		return err
	}

	vars, err := client.GetEnvVars(context.Background(), repo.ID)
	if err != nil {
		return err
	}

	for _, v := range opts.Vars {
		s := strings.Split(v, "=")
		if len(s) != 2 {
			return fmt.Errorf("invalid environment variable syntax. expected KEY=VALUE")
		}

		vars[s[0]] = s[1]
	}

	err = client.SetEnvVars(context.Background(), repo.ID, vars)
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, color.GreenString("Environment variables set"))

	return nil
}
