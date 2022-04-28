package api

import (
	"context"

	graphql "github.com/hasura/go-graphql-client"
)

type Project struct {
	ID graphql.String
}

func GetProject(ctx context.Context, project string) (*Project, error) {
	client := NewGraphQLClient()
	var query struct {
		Project Project `graphql:"project(path: $path)"`
	}
	err := client.Query(ctx, &query, map[string]interface{}{
		"path": graphql.String(project),
	})
	if err != nil {
		return nil, err
	}
	return &query.Project, nil
}

func GetProjectLogs(ctx context.Context, projectID string) ([]string, error) {
	client := NewGraphQLClient()
	var query struct {
		CurrentUser struct {
			Repo struct {
				ProductionDeployment struct {
					Logs []struct {
						Text graphql.String
					}
				}
			} `graphql:"repo(id: $id)"`
		} `graphql:"currentUser()"`
	}
	err := client.Query(ctx, &query, map[string]interface{}{
		"id": graphql.ID(projectID),
	})
	if err != nil {
		return nil, err
	}

	logs := []string{}
	for _, log := range query.CurrentUser.Repo.ProductionDeployment.Logs {
		logs = append(logs, string(log.Text))
	}

	return logs, nil
}
