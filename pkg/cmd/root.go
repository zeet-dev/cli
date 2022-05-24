/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeet-dev/cli/internal/config"
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
	rootCmd.AddCommand(NewDeployCmd(f))
	rootCmd.AddCommand(NewRestartCmd(f))
	rootCmd.AddCommand(NewStatusCmd(f))
	rootCmd.AddCommand(NewEnvSetCmd(f))
	rootCmd.AddCommand(NewEnvGetCmd(f))
	rootCmd.AddCommand(NewConfigSetCmd(f))
	rootCmd.AddCommand(NewJobRunCmd(f))

	rootCmd.AddCommand(NewGenDocsCmd())

	// Set inputs/outputs
	rootCmd.SetErr(&cmdutil.ErrorWriter{Out: f.IOStreams.Out})
	rootCmd.SetIn(f.IOStreams.In)
	rootCmd.SetOut(f.IOStreams.Out)

	// Config & flags
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", filepath.Join(configDir(), defaultConfigName), "Config file")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindEnv("server")
	viper.BindEnv("ws-server")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	return rootCmd
}

// TODO put viper code somewhere else?
func initConfig() {
	viper.SetEnvPrefix("ZEET")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
			err := writeDefaultConfig()
			cobra.CheckErr(err)
		} else {
			cobra.CheckErr(err)
		}
	}

	if viper.GetBool("debug") {
		fmt.Println("Using " + viper.ConfigFileUsed())
	}
}

func writeDefaultConfig() error {
	viper.Set("server", "https://anchor.zeet.co")
	viper.Set("ws-server", "wss://anchor.zeet.co")
	return viper.WriteConfig()
}

func configDir() string {
	p, err := config.ConfigDir()
	cobra.CheckErr(err)
	return p
}
