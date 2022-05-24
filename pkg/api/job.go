package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type Job struct {
	ID    uuid.UUID `copier:"Id"`
	State JobRunState
}

func (c *Client) RunJob(ctx context.Context, projectID uuid.UUID, command string, build bool) (*Job, error) {
	out := &Job{}

	_ = `# @genqlient
		mutation runJob($id: UUID!, $command: String!, $build: Boolean!) {
		  runJob(input: {id: $id, runCommand: $command, build: $build}) {
			state
			id
		  }
		}
	`
	res, err := runJob(ctx, c.GQL, projectID, command, build)
	if err != nil {
		return out, err
	}

	if err := copier.Copy(&out, res.RunJob); err != nil {
		return out, err
	}

	return out, nil
}

func (c *Client) GetJobLogs(ctx context.Context, projectID uuid.UUID, jobID uuid.UUID) ([]LogEntry, error) {
	out := []LogEntry{}

	_ = `# @genqlient
		query getJobLogs($projectID: UUID!, $jobID: UUID!) {
		  project(id: $projectID) {
			repo {
			  jobRun(id: $jobID) {
				logs {
				  entries {
					text
					timestamp
				  }
				}
			  }
			}
		  }
		}
	`
	res, err := getJobLogs(ctx, c.GQL, projectID, jobID)
	if err != nil {
		return out, err
	}

	if err := copier.Copy(&out, res.Project.Repo.JobRun.Logs.Entries); err != nil {
		return out, err
	}

	return out, nil
}

func (c *Client) GetJob(ctx context.Context, projectID uuid.UUID, jobID uuid.UUID) (*Job, error) {
	out := &Job{}

	_ = `# @genqlient
		query getJob($projectID: UUID!, $jobID: UUID!) {
		  project(id: $projectID) {
			repo {
			  jobRun(id: $jobID) {
				id
				state
			  }
			}
		  }
		}
	`
	res, err := getJob(ctx, c.GQL, projectID, jobID)
	if err != nil {
		return out, err
	}

	if err := copier.Copy(&out, res.Project.Repo.JobRun); err != nil {
		return out, err
	}

	return out, nil
}
