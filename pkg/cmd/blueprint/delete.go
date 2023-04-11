package blueprint

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/zeet-dev/cli/pkg/cmdutil"
)

type BlueprintDeleteOptions struct {
	*cmdutil.Factory

	BlueprintID uuid.UUID
}

func NewBlueprintDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command {
		Use: "delete <blueprint_id",
		Short: "Delete a blueprint",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf(color.RedString("Invalid blueprint ID: '%s'\n"), args[0])
				return err
			}

			opts := &BlueprintDeleteOptions{f, id}
			return runBlueprintDelete(opts)
		},
	}
}

func runBlueprintDelete(opts *BlueprintDeleteOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	b, err := client.GetBlueprint(ctx, opts.BlueprintID)
	if err != nil {
		return err
	}

	fmt.Println("Preparing to delete the following blueprint:")
	displayBlueprint(b)

	fmt.Println(color.RedString("WARNING! Deleting a blueprint is irreversible and cannot be undone!"))
	fmt.Printf("Are you sure you want to continue? Enter %s to confirm.\n", color.YellowString(b.Slug))

	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	confirm := strings.TrimSpace(data)

	if confirm == b.Slug {
		fmt.Println("Deleting blueprint...")
		return client.DeleteBlueprint(ctx, opts.BlueprintID)
	}

	fmt.Println("Delete blueprint cancelled.")
	return nil
}

