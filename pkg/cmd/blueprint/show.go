package blueprint

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

type BlueprintShowOptions struct {
	*cmdutil.Factory

	BlueprintID uuid.UUID
}

func NewBlueprintShowCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use: "show <blueprint_id>",
		Short: "View a blueprint",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id, err := uuid.Parse(args[0])

			if err != nil {
				fmt.Printf(color.RedString("Invalid blueprint ID: '%s'\n"), args[0])
				return err
			}

			opts := &BlueprintShowOptions{f, id}

			return runBlueprintShow(opts)
		},
	}
}

func runBlueprintShow(opts *BlueprintShowOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	b, err := client.GetBlueprint(context.Background(), opts.BlueprintID)
	if err != nil {
		return err
	}

	displayBlueprint(b)

	return nil
}

func displayBlueprint(b *api.BlueprintSummary) {
	jsonFormatter := text.NewJSONTransformer("", "  ")
	tw := table.NewWriter()

	tw.AppendRow(table.Row{"ID", b.Id})
	tw.AppendRow(table.Row{"Slug", b.Slug})
	tw.AppendRow(table.Row{"Name", b.DisplayName})
	tw.AppendRow(table.Row{"Description", b.Description})
	tw.AppendRow(table.Row{"Type", b.Type})
	tw.AppendRow(table.Row{"Project Count", b.ProjectCount})

	if b.RichInputSchema != "" {
		tw.AppendRow(table.Row{"Input Schema", jsonFormatter(b.RichInputSchema)})
	}

	for _, t := range b.Tags {
		tw.AppendRow(table.Row{"Tags", t})
	}

	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	tw.SetStyle(table.StyleRounded)

	fmt.Println(tw.Render())
}
