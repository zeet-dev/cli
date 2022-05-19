package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/internal/config"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

type ConfigSetOptions struct {
	IO     *iostreams.IOStreams
	Config func() (config.Config, error)
	Vars   []string
}

func NewConfigSetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ConfigSetOptions{}
	opts.IO = f.IOStreams
	opts.Config = f.Config

	cmd := &cobra.Command{
		Use:     "config:set [name=value]",
		Example: "$ zeet config:set server=https://anchor.zeet.co",
		Short:   "Add or modify a CLI config variable",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Vars = args[0:]

			return runConfigSet(opts)
		},
	}

	return cmd
}

func runConfigSet(opts *ConfigSetOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	for _, v := range opts.Vars {
		s := strings.Split(v, "=")
		if len(s) != 2 {
			return fmt.Errorf("invalid config variable syntax. expected key=value")
		}

		if err := cfg.Set(s[0], s[1]); err != nil {
			return err
		}
	}

	if err := cfg.WriteConfig(); err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, color.GreenString("Config variables set"))

	return nil
}
