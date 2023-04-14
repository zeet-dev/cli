package project

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/display"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

type ProjectListOptions struct{
	*cmdutil.Factory
}

func NewProjectListCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command {
		Use: "list",
		Short: "List projects",
		Args: cobra.ExactArgs(0),
		RunE: func(c *cobra.Command, _ []string) error {
			opts := &ProjectListOptions{f}
			return runProjectList(opts)
		},
	}
}

func runProjectList(opts *ProjectListOptions) error {
	fetch := func(pageInput *api.PageInput) (*string, []table.Row, *api.PageInfo, int, error) {
		ctx := context.Background()
		client, err := opts.ApiClient()
		if err != nil {
			return nil, nil, nil, 0, err
		}

		filter := api.FilterInput{
			Page: *pageInput,
			Filter: api.FilterNode{},
		}

		result, err := client.ListProjectV3s(ctx, filter)
		if err != nil {
			return nil, nil, nil, 0, err
		}

		info := result.PageInfo
		title := "Projects"
		rows := []table.Row{
			table.Row{
				"ID",
				"Name",
			},
		}

		for _, p := range result.Nodes {
			rows = append(rows, table.Row{
				fmt.Sprintf("%s", p.Id),
				p.Name,
			})
		}

		return &title, rows, &info, result.TotalCount, nil
	}

	return display.DisplayPaginatedTable(fetch)
}
