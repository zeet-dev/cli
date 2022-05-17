/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

var defaultConfigName = "config.yaml"

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

	rootCmd.AddCommand(NewGenDocsCmd())

	// Set inputs/outputs
	rootCmd.SetErr(&cmdutil.ErrorWriter{Out: f.IOStreams.Out})
	rootCmd.SetIn(f.IOStreams.In)
	rootCmd.SetOut(f.IOStreams.Out)

	// Config & flags
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config", "c", filepath.Join(configHome(), defaultConfigName), "Config file")
	rootCmd.PersistentFlags().String("server", "https://anchor.zeet.co", "Zeet API Server")
	rootCmd.PersistentFlags().String("ws-server", "wss://anchor.zeet.co", "Zeet Websocket/Subscriptions Server")
	rootCmd.PersistentFlags().BoolP("debug", "v", false, "Enable verbose debug logging")

	rootCmd.PersistentFlags().MarkHidden("server")
	rootCmd.PersistentFlags().MarkHidden("ws-server")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("ws-server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	return rootCmd
}

func initConfig() {
	viper.SetEnvPrefix("ZEET")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	cfgFile := viper.GetString("config")
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
			// No problem, the config file will be created after login
		} else {
			cobra.CheckErr(err)
		}
	}

	if viper.GetBool("debug") {
		fmt.Println("Using " + viper.ConfigFileUsed())
	}
}

func configHome() string {
	if os.Getenv("APP_ENV") == "gen_docs" {
		return "/your/config/dir/zeet"
	}

	cfgDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	ch := filepath.Join(cfgDir, "zeet")
	err = os.MkdirAll(ch, os.ModePerm)
	cobra.CheckErr(err)

	return ch
}
