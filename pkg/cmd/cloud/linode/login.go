package linode

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

type LinodeLoginOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	CloudID uuid.UUID
}

func NewLinodeLoginCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &LinodeLoginOptions{}
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

			return runLinodeLogin(opts)
		},
	}

	return cmd
}

func runLinodeLogin(opts *LinodeLoginOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	cloud, err := client.GetCloudLinode(context.Background(), opts.CloudID)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ofile := path.Join(home, ".zeet", "clouds", cloud.CurrentUser.LinodeAccount.Id.String(), "env.sh")

	if err := os.MkdirAll(filepath.Dir(ofile), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(ofile, []byte(fmt.Sprintf(`#!/bin/sh
export LINODE_CLI_TOKEN=%s
export LINODE_TOKEN=$LINODE_CLI_TOKEN
echo "Linode credentials configured"
`,
		cloud.CurrentUser.LinodeAccount.AccessToken)), 0600); err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, color.GreenString("Linode creds fetched"))
	return nil
}
