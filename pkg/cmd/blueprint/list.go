package blueprint

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jinzhu/copier"
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

	reader := bufio.NewReader(os.Stdin)
	pageInput := api.PageInput{
		First: 10,
		After: "0",
	}

	for ; ; {
		info, err := fetchAndShowPage(client, pageInput)
		if err != nil {
			return err
		}

		prompt := ""
		if info.HasNextPage {
			prompt = prompt + "[n]ext page | "
		}
		if info.HasPreviousPage {
			prompt = prompt + "[p]revious page | "
		}
		if prompt != "" {
			prompt = prompt + "[q]uit"
		} else {
			return nil
		}

		fmt.Println(prompt)

		data, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		action := strings.ToLower(strings.TrimSpace(data))

		switch action {
		case "q", "Q":
			return nil
		case "n", "N":
			if info.HasNextPage {
				pageInput.After = info.EndCursor
			}
		case "p", "P":
			if info.HasPreviousPage {
				// the "Before" field isn't implemented on the server
				// so we'll do it manually
				start, err := strconv.Atoi(info.StartCursor)
				if err != nil {
					return err
				}
				pageInput.After = strconv.Itoa(start - pageInput.First)
			}
		}
	}

	return nil
}

func fetchAndShowPage(client *api.Client, pageInput api.PageInput) (*api.PageInfo, error) {
		result, err := client.ListBlueprints(context.Background(), pageInput)
		if err != nil {
			return nil, err
		}

		blueprints := make([]*api.BlueprintSummary, 0)
		info := result.PageInfo

		for _, n := range result.Nodes {
			bps := &api.BlueprintSummary{}
			if err := copier.Copy(bps, n.BlueprintSummary); err != nil {
				return nil, err
			}

			blueprints = append(blueprints, bps)
		}

		displayBlueprintListTable(blueprints)

		start, err := strconv.Atoi(info.StartCursor)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Showing blueprints %d-%s of %d.\n", start + 1, info.EndCursor, result.TotalCount)

		return &info, nil
}

func displayBlueprintListTable(blueprints []*api.BlueprintSummary) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Description",
		"Slug",
		"Type",
	})

	for _, b := range blueprints {
		tw.AppendRow(table.Row{
			fmt.Sprintf("%s", b.Id),
			b.DisplayName,
			text.WrapSoft(b.Description, 40),
			b.Slug,
			fmt.Sprintf("%s", b.Type),
		})
	}

	tw.SetTitle("Blueprints")
	tw.SetStyle(table.StyleRounded)
	tw.Style().Options.SeparateRows = true
	fmt.Println(tw.Render())
}
