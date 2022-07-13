/*
Copyright © 2022 Zeet, Inc - All Rights Reserved
test
*/

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
)

type DeleteOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	Project string
}

func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &DeleteOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	deleteCmd := &cobra.Command{
		Use:     "delete [project]",
		Short:   "Delete a project",
		Example: "zeet delete zeet-demo/zeet-demo-node-sample",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]

			return runDelete(opts)
		},
	}

	return deleteCmd
}

func runDelete(opts *DeleteOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	path, err := client.ToProjectPath(opts.Project)
	if err != nil {
		return err
	}

	project, err := client.GetProjectByPath(context.Background(), path)
	if err != nil {
		return err
	}

	branch := opts.Project
	if branch == "" {
		branch, err = client.GetProductionBranch(context.Background(), project.ID)
		if err != nil {
			return err
		}
	}

	err = client.DeleteRepo(context.Background(), project.ID, branch)
	if err != nil {
		return err
	}
	return nil
}
