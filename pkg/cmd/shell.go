/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// shellCmd represents the exec command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Get a shell into zeet project",
	RunE: func(cmd *cobra.Command, args []string) error {
		LoginGate()
		shellExec := []string{"sh", "-c", "clear; (bash || zsh || ash || sh)"}
		return execCmd.RunE(cmd, append(args, shellExec...))
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
