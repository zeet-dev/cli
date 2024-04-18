package api

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	Id    uuid.UUID
	Login string
}

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	_ = `# @genqlient
		query getCurrentUser {
		  currentUser {
			id
			login
		  }
		}
	`
	res, err := GetCurrentUserQuery(ctx, c.gql)
	if err != nil {
		return nil, err
	}

	return &User{Id: uuid.MustParse(res.CurrentUser.Id), Login: res.CurrentUser.Login}, nil
}
