/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
	"golang.org/x/term"
)

type LoginOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)
	Config    func() (config.Config, error)

	AccessToken string
}

func NewLoginCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{}
	opts.IO = f.IOStreams
	opts.Config = f.Config
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Zeet. You'll be prompted for a token (from https://zeet.co/account/api) if it's not passed via --token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(opts)
		},
		Annotations: map[string]string{"skipAuthCheck": "true"},
	}

	cmd.Flags().StringVarP(&opts.AccessToken, "token", "t", "", "Your Zeet access token")

	return cmd
}

func runLogin(opts *LoginOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	apiClient, err := opts.ApiClient()
	if err != nil {
		return err
	}

	accessToken := cfg.GetString("auth.access_token")

	if accessToken != "" {
		if user, err := apiClient.GetCurrentUser(context.Background()); err == nil {
			fmt.Fprintln(opts.IO.Out, "You are logged in as: "+user.Login)
			fmt.Fprintf(opts.IO.Out, "Login as a different user? [y/N]: ")

			reader := bufio.NewReader(os.Stdin)
			data, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			confirm := strings.ToLower(strings.TrimSpace(data))
			if !(confirm == "y" || confirm == "yes") {
				return nil
			}
		}
	}

	// If no access token is provided, prompt for one
	// If an access token is provided, use it
	var newToken string
	if opts.AccessToken == "" {
		fmt.Fprint(opts.IO.Out, "Enter Zeet API token (input hidden): ")
		_newToken, err := term.ReadPassword(int(syscall.Stdin))
		newToken = string(_newToken)
		if err != nil {
			return err
		}
		fmt.Fprintln(opts.IO.Out)
	} else {
		newToken = opts.AccessToken
	}

	if err := cfg.Set("auth.access_token", strings.TrimSpace(newToken)); err != nil {
		return err
	}
	if err := cfg.WriteConfig(); err != nil {
		return err
	}

	// Refresh api client to use updated access token
	apiClient, err = opts.ApiClient()
	if err != nil {
		return err
	}

	user, err := apiClient.GetCurrentUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Fprintln(opts.IO.Out, "You are logged in as: "+user.Login)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
