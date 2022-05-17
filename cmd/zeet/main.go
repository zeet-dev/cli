/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package main

import (
	"os"

	"github.com/zeet-dev/cli/internal/build"
	"github.com/zeet-dev/cli/pkg/cmd"
	"github.com/zeet-dev/cli/pkg/cmd/factory"
)

func main() {
	f := factory.New(build.Version)

	rootCmd := cmd.NewRootCmd(f)

	//rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
	//
	//}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
