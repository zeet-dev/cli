package cluster

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

type KubeconfigGetOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	File      string
	ClusterID uuid.UUID
}

func NewKubeconfigGetCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &KubeconfigGetOptions{}
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

func runKubeconfigGet(opts *KubeconfigGetOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	cluster, err := client.GetClusterKubeconfig(context.Background(), opts.ClusterID)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ofile := path.Join(home, ".zeet", "clusters", cluster.ID.String(), "kubeconfig.yaml")

	if err := os.MkdirAll(filepath.Dir(ofile), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(ofile, []byte(cluster.Kubeconfig), 0600); err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, color.GreenString(fmt.Sprintf("Cluster %s kubeconfig fetched", cluster.Name)))
	return nil
}
