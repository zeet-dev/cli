/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

type LogoutOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)
	Config    func() (config.Config, error)

	AccessToken string
	Overwrite   bool
}

func NewLogoutCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LogoutOptions{}
	opts.IO = f.IOStreams
	opts.Config = f.Config
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout to Zeet",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(opts)
		},
	}

	return cmd
}

func runLogout(opts *LogoutOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	if err := cfg.Set("auth.access_token", ""); err != nil {
		return err
	}
	if err := cfg.WriteConfig(); err != nil {
		return err
	}

	return nil
}
