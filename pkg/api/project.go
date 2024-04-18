package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type Project struct {
	ID uuid.UUID `copier:"Id"`
}

func (c *Client) GetProjectByPath(ctx context.Context, project string) (*Project, error) {
	out := &Project{}

	_ = `# @genqlient
		query getProjectByPath($path: String) {
		  repo(path: $path) {
			id
		  }
		}
	`
	res, err := GetProjectByPathQuery(ctx, c.gql, project)
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(out, res.Repo); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetProjectPath(ctx context.Context, id uuid.UUID) (string, error) {
	_ = `# @genqlient
		query getProjectPath($id: UUID!) {
		  repo(id: $id) {
			fullPath
          }
		}
	`
	res, err := GetProjectPathQuery(ctx, c.gql, id)
	if err != nil {
		return "", err
	}

	return res.Repo.FullPath, nil
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
	res, err := GetProductionBranchQuery(ctx, c.gql, projectID.String())
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
	_, err := UpdateBranchMutation(ctx, c.gql, image, deploy, projectID, branchName)
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
	_, err := UpdateProjectMutation(ctx, c.gql, projectID.String(), image)
	if err != nil {
		return err
	}

	return err
}

// ToProjectPath returns the path for a given project UUID or path
func (c *Client) ToProjectPath(input string) (string, error) {
	id, err := uuid.Parse(input)

	if err == nil {
		return c.GetProjectPath(context.Background(), id)
	} else {
		return input, nil
	}

}
