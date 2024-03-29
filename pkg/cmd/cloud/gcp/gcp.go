package gcp

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
	Name     = "GCP"
	NameFull = "Google Cloud Platform"
)

type GCPLoginOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	CloudID uuid.UUID
}

var eval bool

func NewGCPLoginCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &GCPLoginOptions{}
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

			return runGCPLogin(opts)
		},
	}

	cmd.PersistentFlags().BoolVarP(&eval, "eval", "e", false, "eval $(zeet [args])")

	return cmd
}

func runGCPLogin(opts *GCPLoginOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	cloud, err := client.GetCloudGCP(context.Background(), opts.CloudID)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ofile := path.Join(home, ".zeet", "clouds", cloud.CurrentUser.GcpAccount.Id.String(), "env.sh")

	if err := os.MkdirAll(filepath.Dir(ofile), os.ModePerm); err != nil {
		return err
	}

	sfile := path.Join(home, ".zeet", "clouds", cloud.CurrentUser.GcpAccount.Id.String(), "credentials.json")
	if err := os.WriteFile(sfile, []byte(cloud.CurrentUser.GcpAccount.Credentials), 0600); err != nil {
		return err
	}

	if err := os.WriteFile(ofile, []byte(fmt.Sprintf(`#!/bin/sh
export GOOGLE_APPLICATION_CREDENTIALS=%s
gcloud auth activate-service-account --key-file $GOOGLE_APPLICATION_CREDENTIALS
gcloud config set project %s
echo "GCP credentials configured"
`,
		sfile, cloud.CurrentUser.GcpAccount.ProjectID)), 0600); err != nil {
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
