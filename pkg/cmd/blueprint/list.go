package blueprint

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/display"
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
	fetch := func(pageInput *api.PageInput) (*string, []table.Row, *api.PageInfo, int, error) {
		ctx := context.Background()

		client, err := opts.ApiClient()
		if err != nil {
			return nil, nil, nil, 0, err
		}

		result, err := client.ListBlueprints(ctx, *pageInput)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		info := result.PageInfo
		title := "Blueprints"
		rows := []table.Row{
			table.Row{
				"ID",
				"Name",
				"Description",
				"Slug",
				"Type",
			},
		}

		for _, b := range result.Nodes {
			rows = append(rows, table.Row{
				fmt.Sprintf("%s", b.Id),
				b.DisplayName,
				text.WrapSoft(b.Description, 40),
				b.Slug,
				fmt.Sprintf("%s", b.Type),
			})
		}

		return &title, rows, &info, result.TotalCount, nil
	}

	return display.DisplayPaginatedTable(fetch)
}
