/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/internal/build"
	"github.com/zeet-dev/cli/pkg/cmd"
	"github.com/zeet-dev/cli/pkg/cmd/factory"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func main() {
	f := factory.New(build.Version)
	cfg, err := f.Config()
	if err != nil {
		fmt.Fprintf(f.IOStreams.ErrOut, "failed to read configuration: %s\n", err)
		os.Exit(1)
	}

	rootCmd := cmd.NewRootCmd(f)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmdutil.IsAuthCheckEnabled(cmd) && !cmdutil.CheckAuth(cfg) {
			return fmt.Errorf("not logged in (hint: run 'zeet login')")
		}

		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
