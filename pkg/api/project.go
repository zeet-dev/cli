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

func (c *Client) UpdateBranch(ctx context.Context, projectID uuid.UUID, image string, branchName string, deploy bool) error {
	_ = `# @genqlient
		mutation updateBranch($image: String!, $deploy: Boolean, $projectID: UUID!, $branchName: String!) {
		  updateBranch(input: {image: $image, deploy: $deploy, repoID: $projectID, name: $branchName}) {
			id
		  }
		}
	`
	_, err := updateBranch(ctx, c.GQL, image, deploy, projectID, branchName)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) UpdateProject(ctx context.Context, projectID uuid.UUID, image string) error {
	_ = `# @genqlient
		mutation updateProject($projectID: ID!, $image: String!) {
		  updateProject(input: {id: $projectID, dockerImage: $image}) {
			id
		  }
		}
	`
	_, err := updateProject(ctx, c.GQL, projectID, image)
	if err != nil {
		return err
	}

	return err
}
