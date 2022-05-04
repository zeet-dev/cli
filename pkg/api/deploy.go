package api

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/exp/slices"
)

type Deployment struct {
	ID              uuid.UUID        `copier:"Id"`
	Status          DeploymentStatus `copier:"Status"`
	Branch          string
	Endpoints       []string
	PrivateEndpoint string
}

type LogEntry struct {
	Text      string
	Timestamp time.Time
}

func (c *Client) BuildProject(ctx context.Context, projectID uuid.UUID, branch string, noCache bool) (*Deployment, error) {
	_ = `# @genqlient
		mutation buildRepo($id: ID!, $branch: String, $noCache: Boolean) {
		  buildRepo(id: $id, branch: $branch, noCache: $noCache) {
			deployments {
			  id
			  status
              branch
			  endpoints
			  privateEndpoint
			}
		  }
		}
	`
	_ = `# @genqlient
		mutation buildRepoDefaultBranch($id: ID!, $noCache: Boolean) {
		  buildRepo(id: $id, noCache: $noCache) {
			deployments {
			  id
			  status
              branch
			  endpoints
			  privateEndpoint
			}
		  }
		}
	`

	out := &Deployment{}

	if branch == "" {
		res, err := buildRepoDefaultBranch(ctx, c.GQL, projectID, noCache)
		if err != nil {
			return nil, err
		}
		if err := copier.Copy(out, res.BuildRepo.Deployments[0]); err != nil {
			return nil, err
		}
	} else {
		res, err := buildRepo(ctx, c.GQL, projectID, branch, noCache)
		if err != nil {
			return nil, err
		}
		if err := copier.Copy(out, res.BuildRepo.Deployments[0]); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (c *Client) GetBuildLogs(ctx context.Context, deployID uuid.UUID) (out []LogEntry, err error) {
	_ = `# @genqlient
		query getBuildLogs($id: ID!) {
		  currentUser {
			deployment(id: $id) {
			  build {
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

	res, err := getBuildLogs(ctx, c.GQL, deployID)
	if err != nil {
		return
	}

	err = copier.Copy(&out, res.CurrentUser.Deployment.Build.Logs.Entries)
	return
}

func (c *Client) GetDeploymentLogs(ctx context.Context, deployID uuid.UUID) (out []LogEntry, err error) {
	_ = `# @genqlient
		query getDeploymentLogs($id: ID!) {
		  currentUser {
			deployment(id: $id) {
				deployStep {
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

	res, err := getDeploymentLogs(ctx, c.GQL, deployID)
	if err != nil {
		return
	}

	err = copier.Copy(&out, res.CurrentUser.Deployment.DeployStep.Logs.Entries)
	return
}

func (c *Client) GetDeployment(ctx context.Context, deploymentID uuid.UUID) (*Deployment, error) {
	out := &Deployment{}

	_ = `# @genqlient
		query getDeploymentInfo($id: ID!) {
		  currentUser {
			deployment(id: $id) {
			  id
			  status
			  endpoints
			  privateEndpoint
			}
		  }
		}
	`
	res, err := getDeploymentInfo(ctx, c.GQL, deploymentID)
	if err != nil {
		return nil, err
	}

	err = copier.Copy(out, res.CurrentUser.Deployment)
	return out, err
}

func IsDeployInProgress(status DeploymentStatus) bool {
	inProgress := []DeploymentStatus{DeploymentStatusDeployPending,
		DeploymentStatusDeployInProgress,
		DeploymentStatusReleaseInProgress,
	}
	return slices.Contains(inProgress, status)
}

func IsDeployFailed(status DeploymentStatus) bool {
	failed := []DeploymentStatus{
		DeploymentStatusDeployFailed,
		DeploymentStatusDeployStopped,
		DeploymentStatusDeployCrashing,
	}

	return slices.Contains(failed, status)
}

func IsDeploySuccess(status DeploymentStatus) bool {
	failed := []DeploymentStatus{
		DeploymentStatusDeploySucceeded,
		DeploymentStatusDeployHealhty,
	}

	return slices.Contains(failed, status)
}

func IsBuildInProgress(status DeploymentStatus) bool {
	inProgress := []DeploymentStatus{
		DeploymentStatusBuildPending,
		DeploymentStatusBuildInProgress,
	}
	return slices.Contains(inProgress, status)
}

func IsBuildFailed(status DeploymentStatus) bool {
	failed := []DeploymentStatus{
		DeploymentStatusBuildFailed,
		DeploymentStatusBuildAborted,
	}

	return slices.Contains(failed, status)
}

func IsBuildSuccess(status DeploymentStatus) bool {
	failed := []DeploymentStatus{
		DeploymentStatusBuildSucceeded,
		DeploymentStatusDeployPending,
	}

	return slices.Contains(failed, status)
}
