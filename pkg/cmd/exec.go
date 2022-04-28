/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute command in Zeet project",
	Run: func(cmd *cobra.Command, args []string) {
		LoginGate()
		fmt.Println("exec called")
	},
}
