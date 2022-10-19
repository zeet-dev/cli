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

const (
	Name     = "Cluster"
	NameFull = "Kubernetes Cluster"
)

type ClusterLoginOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	ClusterID uuid.UUID
}

var eval bool

func NewClusterLoginCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &ClusterLoginOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:  "login <cluster id>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := uuid.Parse(args[0])
			if err != nil {
				return errors.New("invalid cluster ID format")
			}
			opts.ClusterID = id

			return runClusterLogin(opts)
		},
	}

	cmd.PersistentFlags().BoolVarP(&eval, "eval", "e", false, "eval $(zeet [args])")

	return cmd
}

func runClusterLogin(opts *ClusterLoginOptions) error {
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

	if eval {
		fmt.Fprintln(opts.IO.ErrOut, color.GreenString(fmt.Sprintf("%s %s creds fetched", Name, cluster.Name)))
		fmt.Fprintf(opts.IO.Out, "export KUBECONFIG=%s", ofile)
	} else {
		fmt.Fprintln(opts.IO.Out, color.GreenString(fmt.Sprintf("%s %s creds fetched", Name, cluster.Name)))
	}

	return nil
}
