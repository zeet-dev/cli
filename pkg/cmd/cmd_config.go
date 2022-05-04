package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/config"
)

type CmdConfig struct {
	cmd    *cobra.Command
	cfg    config.Provider
	args   []string
	ctx    context.Context
	client *api.Client
}

func NewCmdConfig(cmd *cobra.Command, args []string) *CmdConfig {
	cfg := &config.Live{}
	client := cfg.GetAPIClient(cfg.GetString("server"), cfg.GetString("auth.access_token"))

	return &CmdConfig{cfg: cfg, cmd: cmd, args: args, ctx: context.Background(), client: client}
}

func withCmdConfig(runner func(cmdConfig *CmdConfig) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmdConfig := NewCmdConfig(cmd, args)

		return runner(cmdConfig)
	}
}
