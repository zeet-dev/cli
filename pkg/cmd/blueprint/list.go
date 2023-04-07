package blueprint

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

type BlueprintListOptions struct{
	*cmdutil.Factory
}

func NewBlueprintListCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Short: "List blueprints",
		Args: cobra.ExactArgs(0),
		RunE: func(c *cobra.Command, _ []string) error {
			opts := &BlueprintListOptions{f}
			return runBlueprintList(opts)
		},
	}
}

func runBlueprintList(opts *BlueprintListOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	bps, err := client.ListBlueprints(context.Background())
	if err != nil {
		return err
	}

	displayBlueprintListTable(bps)

	return nil
}

func displayBlueprintListTable(blueprints []*api.BlueprintSummary) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Description",
		"Type",
	})

	for _, b := range blueprints {
		tw.AppendRow(table.Row{
			fmt.Sprintf("%s", b.Id),
			b.DisplayName,
			b.Description,
			fmt.Sprintf("%s", b.Type),
		})
	}

	tw.SetTitle("Blueprints")
	fmt.Println(tw.Render())
}
