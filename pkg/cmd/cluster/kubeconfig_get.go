package cluster

import (
	"errors"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func NewKubeconfigGetCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &ClusterLoginOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:   "kubeconfig:get <cluster id> [kubeconfig location]",
		Short: "Download a kubeconfig from Zeet",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := uuid.Parse(args[0])
			if err != nil {
				return errors.New("invalid cluster ID format")
			}
			opts.ClusterID = id

			return runKubeconfigGet(opts)
		},
	}

	return cmd
}

func runKubeconfigGet(opts *ClusterLoginOptions) error {
	return runClusterLogin(opts)
}
