package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

const (
	Name     = "DO"
	NameFull = "Digital Ocean"
)

type DOLoginOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	CloudID uuid.UUID
}

var eval bool

func NewDOLoginCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &DOLoginOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:  "login <cloud id>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := uuid.Parse(args[0])
			if err != nil {
				return errors.New("invalid cloud ID format")
			}
			opts.CloudID = id

			return runDOLogin(opts)
		},
	}

	cmd.PersistentFlags().BoolVarP(&eval, "eval", "e", false, "eval $(zeet [args])")

	return cmd
}

func runDOLogin(opts *DOLoginOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	cloud, err := client.GetCloudDigitalOcean(context.Background(), opts.CloudID)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ofile := path.Join(home, ".zeet", "clouds", cloud.CurrentUser.DoAccount.Id.String(), "env.sh")

	if err := os.MkdirAll(filepath.Dir(ofile), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(ofile, []byte(fmt.Sprintf(`#!/bin/sh
export DIGITALOCEAN_ACCESS_TOKEN=%s
echo "DO credentials configured"
`,
		cloud.CurrentUser.DoAccount.AccessToken)), 0600); err != nil {
		return err
	}

	if eval {
		fmt.Fprintln(opts.IO.ErrOut, color.GreenString(fmt.Sprintf("%s creds fetched", Name)))
		fmt.Fprintf(opts.IO.Out, "source %s", ofile)
	} else {
		fmt.Fprintln(opts.IO.Out, color.GreenString(fmt.Sprintf("%s creds fetched", Name)))
	}

	return nil
}
