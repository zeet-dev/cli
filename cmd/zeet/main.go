/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/internal/update"
	"github.com/zeet-dev/cli/pkg/cmd"
	"github.com/zeet-dev/cli/pkg/cmd/factory"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/utils"
)

func main() {
	version := "v0.1.0"

	f := factory.New(version)
	stdout := f.IOStreams.Out
	stderr := f.IOStreams.ErrOut

	cfg, err := f.Config()
	if err != nil {
		fmt.Fprintf(stderr, "failed to read configuration: %s\n", err)
		os.Exit(1)
	}

	updateMessageChan := make(chan *update.ReleaseInfo)
	go func() {
		// TODO improve error handling?
		state, err := config.NewState()
		if err != nil {
			updateMessageChan <- nil
		}

		rel, _ := checkForUpdate(state, &http.Client{Timeout: time.Second * 3}, version)
		updateMessageChan <- rel
	}()

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

	newRelease := <-updateMessageChan
	if newRelease != nil {
		isHomebrew := isFromHomebrew()

		fmt.Fprintf(stdout, "\n%s %s -> %s",
			color.YellowString("A new release of zeet is available:"),
			color.YellowString(version),
			color.GreenString(newRelease.TagName),
		)

		if isHomebrew {
			fmt.Fprintf(stdout, color.YellowString("\nTo upgrade, run: %s\n", "brew update && brew upgrade zeet"))
		}

		fmt.Fprintf(stderr, "\n%s\n",
			color.YellowString(newRelease.URL))
	}
}

func checkForUpdate(state config.Config, client *http.Client, currentVersion string) (*update.ReleaseInfo, error) {
	if utils.IsCI() || !update.ShouldCheck(state) {
		return nil, nil
	}

	return update.Check(client, state, currentVersion)
}

func isFromHomebrew() bool {
	exe, err := exec.LookPath("brew")
	if err != nil {
		return false
	}
	_, err = exec.Command(exe, "list", "zeet").Output()
	return err == nil
}
