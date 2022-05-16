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
