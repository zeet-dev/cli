/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeet-dev/cli/pkg/utils"
)

var defaultConfigName = "config.yaml"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "zeet",
	Short:        "Zeet CLI",
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetErr(&utils.ErrorWriter{})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("config", "c", filepath.Join(configHome(), defaultConfigName), "Config file")
	rootCmd.PersistentFlags().StringP("server", "s", "https://anchor.zeet.co", "Zeet API Server")
	rootCmd.PersistentFlags().String("ws-server", "wss://anchor.zeet.co", "Zeet Websocket/Subscriptions Server")
	rootCmd.PersistentFlags().BoolP("debug", "v", false, "Enable verbose debug logging")

	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("ws-server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func initConfig() {
	viper.SetEnvPrefix("ZEET")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	cfgFile := viper.GetString("config")
	viper.SetConfigFile(cfgFile)

	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	if viper.GetBool("debug") {
		fmt.Println("Using " + viper.ConfigFileUsed())
	}
}

func configHome() string {
	cfgDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	ch := filepath.Join(cfgDir, "zeet")
	err = os.MkdirAll(ch, os.ModePerm)
	cobra.CheckErr(err)

	return ch
}
