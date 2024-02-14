package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type Repo struct {
	ID uuid.UUID `copier:"Id"`
}

func (c *Client) DeleteRepo(ctx context.Context, repoID uuid.UUID) error {
	_ = `# @genqlient
		mutation delete($id: ID!) {
		  deleteRepo(id: $id)
		}
	`

	_, err := DeleteMutation(ctx, c.gql, repoID)
	return err
}

func (c *Client) SetEnvVars(ctx context.Context, repoID uuid.UUID, vars map[string]string) error {
	var inp []EnvVarInput
	for k, v := range vars {
		inp = append(inp, EnvVarInput{
			Name:  k,
			Value: v,
		})
	}

	_ = `# @genqlient
		mutation setEnvVars($id: ID!, $envs: [EnvVarInput!]!) {
		  setRepoEnvs(input: {id: $id, envs: $envs}) {
			envs {
			  id
			}
		  }
		}
	`
	_, err := SetEnvVarsMutation(ctx, c.gql, repoID, inp)
	return err
}

func (c *Client) GetEnvVars(ctx context.Context, repoID uuid.UUID) (map[string]string, error) {
	_ = `# @genqlient
		query getEnvVars($id: ID!) {
		  currentUser {
			repo(id: $id) {
			  envs {
				name
				value
			  }
			}
		  }
		}
	`
	res, err := GetEnvVarsQuery(ctx, c.gql, repoID)
	if err != nil {
		return nil, err
	}

	out := map[string]string{}
	for _, env := range res.CurrentUser.Repo.Envs {
		out[env.Name] = env.Value
	}

	return out, nil
}

func (c *Client) GetRepo(ctx context.Context, path string) (*Repo, error) {
	out := &Repo{}

	_ = `# @genqlient
		query getRepo($path: String!) {
		  repo(path: $path) {
			id
		  }
		}
	`
	res, err := GetRepoQuery(ctx, c.gql, path)
	if err := copier.Copy(out, res.Repo); err != nil {
		return nil, err
	}

	return out, err
}
