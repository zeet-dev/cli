package utils

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeet-dev/cli/pkg/api"
)

func GetProjectByPathOrUUID(client *api.Client, project string) (*api.Project, error) {
	id, err := uuid.Parse(project)
	if err == nil {
		return client.GetProjectById(context.Background(), id)
	} else {
		return client.GetProjectByPath(context.Background(), project)
	}
}

// ToProjectPath returns the path for a given project UUID or path
func ToProjectPath(client *api.Client, input string) (string, error) {
	id, err := uuid.Parse(input)

	if err == nil {
		return client.GetProjectPath(context.Background(), id)
	} else {
		return input, nil
	}

}
