/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/pkg/cmd/cloud"
	"github.com/zeet-dev/cli/pkg/cmd/cluster"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

var defaultConfigName = "config.yaml"
var configPath string

// NewRootCmd creates a cobra.Command and adds subcommands to it
// It's called by main.go and passed a cmdutil.Factory
func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "zeet",
		Short:        "Zeet CLI",
		SilenceUsage: true,
	}

	// Commands
	rootCmd.AddCommand(NewLoginCmd(f))
	rootCmd.AddCommand(NewLogsCmd(f))

	// Project Commands
	rootCmd.AddCommand(NewDeployCmd(f))
	rootCmd.AddCommand(NewRestartCmd(f))
	rootCmd.AddCommand(NewStatusCmd(f))
	rootCmd.AddCommand(NewEnvSetCmd(f))
	rootCmd.AddCommand(NewEnvGetCmd(f))
	rootCmd.AddCommand(NewConfigSetCmd(f))
	rootCmd.AddCommand(NewJobRunCmd(f))
	rootCmd.AddCommand(NewDeleteCmd(f))

	// Cloud Commands
	cloud.InitCloudCmds(f, rootCmd)
	rootCmd.AddCommand(cluster.NewClusterCmd(f))

	// Set inputs/outputs
	rootCmd.SetErr(&cmdutil.ErrorWriter{Out: f.IOStreams.Out})
	rootCmd.SetIn(f.IOStreams.In)
	rootCmd.SetOut(f.IOStreams.Out)

	// Config & flags
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", filepath.Join(configDir(), defaultConfigName), "Config file")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	viper.BindEnv("api-url")
	viper.SetDefault("api-url", "https://anchor.zeet.co")

	viper.BindEnv("ws-url")
	viper.SetDefault("ws-url", "wss://anchor.zeet.co")

	viper.BindEnv("auth.access_token", "ZEET_TOKEN")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	return rootCmd
}

// TODO put viper code somewhere else?
func initConfig() {
	viper.SetEnvPrefix("ZEET")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
		} else {
			cobra.CheckErr(err)
		}
	}

	if viper.GetBool("debug") {
		fmt.Println("Using " + viper.ConfigFileUsed())
	}
}

func configDir() string {
	p, err := config.ConfigDir()
	cobra.CheckErr(err)
	return p
}
