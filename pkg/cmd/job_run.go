/*
Copyright © 2022 Zeet, Inc - All Rights Reserved

*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"github.com/zeet-dev/cli/pkg/iostreams"
	"github.com/zeet-dev/cli/pkg/utils"
)

type JobRunOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.Client, error)

	Project string
	Command string
	Build   bool
	Follow  bool
}

func NewJobRunCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &JobRunOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:   "job:run [project] [command]",
		Short: "Executes a command on a project, in a new instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Project = args[0]
			opts.Command = strings.Join(args[1:], " ")

			return runJobRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Build, "build", "b", false, "Trigger build (true) or use latest image (false)")
	cmd.Flags().BoolVarP(&opts.Follow, "follow", "f", false, "Run until the job is complete")

	cmd.MarkFlagRequired("cmd")

	return cmd
}

func runJobRun(opts *JobRunOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	path, err := utils.ToProjectPath(client, opts.Project)
	if err != nil {
		return err
	}
	project, err := client.GetProjectByPath(context.Background(), path)
	if err != nil {
		return err
	}

	job, err := client.RunJob(context.Background(), project.ID, opts.Command, opts.Build)
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, "Starting job...")
	fmt.Fprintf(opts.IO.Out, "Dashboard: %s\n\n", fmt.Sprintf("https://zeet.co/repo/%s/jobs/%s", project.ID, job.ID))

	if opts.Follow {
		return nil
	}

	jobFinished := false
	logsPrinted := false
	for !jobFinished {
		job, err = client.GetJob(context.Background(), project.ID, job.ID)
		if err != nil {
			return err
		}

		// TODO improve logic...
		switch job.State {
		case api.JobRunStateJobRunStarting:
			break
		case api.JobRunStateJobRunRunning:
			if !logsPrinted {
				if err := pollJobLogs(client, project, job, opts.IO.Out); err != nil {
					return err
				}
			}
			logsPrinted = true
			break
		case api.JobRunStateJobRunSucceeded:
			jobFinished = true
			if !logsPrinted {
				if err := printJobLogs(client, project, job, opts.IO.Out); err != nil {
					return err
				}
			}
			logsPrinted = true
			fmt.Fprintln(opts.IO.Out, color.GreenString("Job ran successfully"))
			break
		case api.JobRunStateJobRunFailed:
			jobFinished = true
			if !logsPrinted {
				if err := printJobLogs(client, project, job, opts.IO.Out); err != nil {
					return err
				}
			}
			logsPrinted = true
			fmt.Fprintln(opts.IO.Out, color.RedString("Job failed"))
			break
		}
	}

	return nil
}

func pollJobLogs(client *api.Client, project *api.Project, job *api.Job, out io.Writer) error {
	getLogs := func() ([]api.LogEntry, error) {
		return client.GetJobLogs(context.Background(), project.ID, job.ID)
	}
	getState := func() (api.JobRunState, error) {
		job, err := client.GetJob(context.Background(), project.ID, job.ID)
		if err != nil {
			return "", err
		}

		return job.State, nil
	}

	return utils.PollLogs(getLogs, getState, out)
}

func printJobLogs(client *api.Client, project *api.Project, job *api.Job, out io.Writer) error {
	logs, err := client.GetJobLogs(context.Background(), project.ID, job.ID)
	if err != nil {
		return err
	}

	for _, l := range logs {
		fmt.Fprintln(out, l.Text)
	}

	return nil
}
