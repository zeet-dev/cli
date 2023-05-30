package project

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

type ProjectShowOptions struct {
	*cmdutil.Factory

	ProjectID uuid.UUID
}

func NewProjectShowCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use: "show <project_id>",
		Short: "View a project",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf(color.RedString("Invalid project ID: '%s'\n"), args[0])
				return err
			}

			opts := &ProjectShowOptions{f, id}
			return runProjectShow(opts)
		},
	}
}

func runProjectShow(opts *ProjectShowOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	p, err := client.GetProjectV3(context.Background(), opts.ProjectID)
	if err != nil {
		return err
	}

	displayProject(p)

	return nil
}

func displayProject(p *api.ProjectV3AdapterSummary) {
	tw := table.NewWriter()

	tw.AppendRow(table.Row{"ID", p.Id})
	tw.AppendRow(table.Row{"Name", p.Name})
	tw.AppendRow(table.Row{"ProjectID", p.GetProjectV3().Id})
	tw.AppendRow(table.Row{"ProjectName", p.GetProjectV3().Name})
	tw.AppendRow(table.Row{"RepoID", p.GetRepo().Id})
	tw.AppendRow(table.Row{"RepoName", p.GetRepo().Name})

	tw.SetStyle(table.StyleRounded)

	fmt.Println(tw.Render())
}
