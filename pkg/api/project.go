package api

import (
	"context"

	"github.com/google/uuid"
)

type Project struct {
	ID uuid.UUID
}

func (c *Client) GetProjectByPath(ctx context.Context, project string) (*Project, error) {
	_ = `# @genqlient
		query getProjectByPath($path: String) {
		  project(path: $path) {
			id
		  }
		}
	`
	res, err := getProjectByPath(ctx, c.GQL, project)
	if err != nil {
		return nil, err
	}
	return &Project{res.Project.Id}, nil
}

func (c *Client) GetProjectLogs(ctx context.Context, projectID uuid.UUID) ([]string, error) {
	_ = `# @genqlient
		query getProjectLogs($path: ID!) {
		  currentUser {
			repo(id: $path) {
			  productionDeployment {
				logs {
				  text
				}
			  }
			}
		  }
		}
	`
	res, err := getProjectLogs(ctx, c.GQL, projectID)
	if err != nil {
		return nil, err
	}

	var logs []string
	for _, log := range res.CurrentUser.Repo.ProductionDeployment.Logs {
		logs = append(logs, log.Text)
	}

	return logs, nil
}

func (c *Client) GetProductionBranch(ctx context.Context, projectID uuid.UUID) (string, error) {
	_ = `# @genqlient
		query getProductionBranch($repoId: ID!) {
		  currentUser {
			repo(id: $repoId) {
			  id
			  productionBranchV2 {
				name
			  }
			}
		  }
		}
	`
	res, err := getProductionBranch(ctx, c.GQL, projectID)
	if err != nil {
		return "", err
	}

	return res.CurrentUser.Repo.ProductionBranchV2.Name, nil
}
