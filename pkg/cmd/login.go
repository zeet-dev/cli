/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"bufio"
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
	RunE: withCmdConfig(func(c *CmdConfig) error {
		c.cmd.SilenceUsage = true

		return Login(c)
	}),
}

func Login(c *CmdConfig) error {
	accessToken := c.cfg.GetString("auth.access_token")

	if accessToken != "" {
		if user, err := c.client.GetCurrentUser(c.ctx); err == nil {
			fmt.Println("You are logged in as: " + user.Login)
			fmt.Print("Login as a different user? [y/N]: ")

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

	fmt.Print("Enter Zeet API token (input hidden): ")
	newToken, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println()

	if err := c.cfg.Set("auth.access_token", strings.TrimSpace(string(newToken))); err != nil {
		return err
	}
	// We'll need to recreate the client so that it uses the updated access token
	c.client = api.New(c.cfg.GetString("server"), c.cfg.GetString("auth.access_token"))

	user, err := c.client.GetCurrentUser(c.ctx)
	if err != nil {
		return err
	}
	fmt.Println("You are logged in as: " + user.Login)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// checkLoginAndRun runs runner if the user is logged in, and returns an error if not
func checkLoginAndRun[O any](c *CmdConfig, runner func(c *CmdConfig, opts O) error, opts O) error {
	accessToken := c.cfg.GetString("auth.access_token")
	if accessToken == "" {
		return fmt.Errorf("not logged in (hint: run 'zeet login')")
	} else {
		return runner(c, opts)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
