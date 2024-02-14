package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/zeet-dev/cli/pkg/utils"
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
	res, err := RunJobMutation(ctx, c.gql, projectID, command, build)
	if err != nil {
		return out, err
	}

	if err := copier.Copy(&out, res.RunJob); err != nil {
		return out, err
	}

	return out, nil
}

func (c *Client) GetJobLogs(ctx context.Context, repoID uuid.UUID, jobID uuid.UUID) ([]LogEntry, error) {
	out := []LogEntry{}

	_ = `# @genqlient
		query getJobLogs($repoID: UUID!, $jobID: UUID!) {
		  repo(id: $repoID) {
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
	`
	res, err := GetJobLogsQuery(ctx, c.gql, repoID, jobID)
	if err != nil {
		return out, err
	}

	if err := copier.Copy(&out, res.Repo.JobRun.Logs.Entries); err != nil {
		return out, err
	}

	return out, nil
}

func (c *Client) GetJob(ctx context.Context, projectID uuid.UUID, jobID uuid.UUID) (*Job, error) {
	out := &Job{}

	_ = `# @genqlient
		query getJob($projectID: UUID!, $jobID: UUID!) {
			repo(id: $projectID) {
			  jobRun(id: $jobID) {
				id
				state
			  }
			}
		}
	`
	res, err := GetJobQuery(ctx, c.gql, projectID, jobID)
	if err != nil {
		return out, err
	}

	if err := copier.Copy(&out, res.Repo.JobRun); err != nil {
		return out, err
	}

	return out, nil
}

func IsJobInProgress(state JobRunState) bool {
	ok := []JobRunState{JobRunStateJobRunStarting, JobRunStateJobRunRunning}
	return utils.SliceContains(ok, state)
}
