/*
Copyright © 2021 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "zeet",
	Short:        "Zeet CLI",
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/zeet/config.yaml)")
	rootCmd.PersistentFlags().StringP("server", "s", "https://anchor.zeet.co", "Zeet API Server")
	rootCmd.PersistentFlags().BoolP("debug", "v", false, "Enable verbose debug logging")

	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFileEnv := os.Getenv("ZEET_CONFIG")
	if cfgFileEnv != "" {
		cfgFile = cfgFileEnv
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".config/zeet" (without extension).
		viper.AddConfigPath(path.Join(home, ".config", "zeet"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("ZEET")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		fmt.Fprintln(os.Stderr, viper.AllSettings())
	}
}
