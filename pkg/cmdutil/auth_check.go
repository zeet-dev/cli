package cmdutil

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/internal/config"
)

func CheckAuth(c config.Config) bool {
	return c.GetString("auth.access_token") != ""
}

func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
		return false
	}

	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["skipAuthCheck"] == "true" {
			return false
		}
	}

	return true
}
