package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
	"gopkg.in/yaml.v2"
)

type ClusterListOptions struct {
	IO           *iostreams.IOStreams
	ApiClient    func() (*api.Client, error)
	OutputFormat string

	Team string
}

func NewClusterListCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &ClusterListOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use: "list [team]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Team = args[0]
			}

			return runClusterList(opts)
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.OutputFormat, "output", "o", "text", "format: text|json|yaml|table")

	return cmd
}

func runClusterList(opts *ClusterListOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	clusters, err := client.ListClusters(context.Background(), opts.Team)
	if err != nil {
		return err
	}

	switch opts.OutputFormat {
	case "text":
		for _, cluster := range clusters {
			fmt.Fprintln(opts.IO.Out, cluster.ID, cluster.Name, cluster.CloudProvider, cluster.ClusterProvider, cluster.Region, cluster.Conencted)
		}
	case "json":
		output, err := json.MarshalIndent(clusters, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Fprintln(opts.IO.Out, string(output))
	case "yaml":
		output, err := yaml.Marshal(clusters)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Fprintln(opts.IO.Out, string(output))
	case "table":
		t := table.NewWriter()
		t.SetOutputMirror(opts.IO.Out)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{
			text.FgGreen.Sprint("ID"),
			text.FgGreen.Sprint("Name"),
			text.FgGreen.Sprint("Cloud Provider"),
			text.FgGreen.Sprint("Cluster Provider"),
			text.FgGreen.Sprint("Region"),
			text.FgGreen.Sprint("Connected"),
		})

		for _, cluster := range clusters {
			t.AppendRows([]table.Row{
				{cluster.ID.String(), cluster.Name, cluster.CloudProvider, cluster.ClusterProvider, cluster.Region, cluster.Conencted},
			})
		}

		t.Render()
	}

	return nil
}
