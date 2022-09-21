package aws

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

type AWSLoginOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	CloudID uuid.UUID
}

func NewAWSLoginCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &AWSLoginOptions{}
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

			return runAWSLogin(opts)
		},
	}

	return cmd
}

func runAWSLogin(opts *AWSLoginOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	cloud, err := client.GetCloudAWS(context.Background(), opts.CloudID)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ofile := path.Join(home, ".zeet", "clouds", cloud.CurrentUser.AwsAccount.Id.String(), "env.sh")

	if err := os.MkdirAll(filepath.Dir(ofile), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(ofile, []byte(fmt.Sprintf(`#!/bin/sh
ASSUMED_SESSION=$(aws sts assume-role \
	--role-arn %s \
	--external-id %s \
	--role-session-name zeet \
	--output text \
	--query='Credentials.[
		join(`+"`=`, [`AWS_ACCESS_KEY_ID`, AccessKeyId]),"+`
		join(`+"`=`, [`AWS_SECRET_ACCESS_KEY`, SecretAccessKey]),"+`
		join(`+"`=`, [`AWS_SESSION_TOKEN`, SessionToken])"+`
	]')
if [ $? -eq 0 ]; then
   eval "export $ASSUMED_SESSION"
   echo "AWS credentials configured"
else
   echo "Failed to assume role"
fi
`,
		cloud.CurrentUser.AwsAccount.RoleARN,
		cloud.CurrentUser.AwsAccount.ExternalID)), 0600); err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, color.GreenString("AWS creds fetched"))
	return nil
}
