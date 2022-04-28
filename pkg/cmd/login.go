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
		tokenb, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		viper.Set("auth.access_token", string(tokenb))
		if err := api.ShowCurrentUser(ctx); err != nil {
			return err
		}

		if err := viper.WriteConfig(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
