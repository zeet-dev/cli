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
	"github.com/zeet-dev/cli/pkg/api"
	"golang.org/x/term"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Zeet",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		ctx := context.Background()

		return Login(ctx)
	},
}

func Login(ctx context.Context) error {
	token := viper.GetString("auth.access_token")
	if token != "" {
		if err := api.ShowCurrentUser(ctx); err == nil {
			fmt.Print("Login as another user? [y/N]: ")

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

	fmt.Print("Enter Zeet API token: ")
	newToken, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println()

	viper.Set("auth.access_token", strings.TrimSpace(string(newToken)))
	if err := api.ShowCurrentUser(ctx); err != nil {
		return err
	}

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func LoginGate() error {
	token := viper.GetString("auth.access_token")
	if token == "" {
		return Login(context.Background())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
