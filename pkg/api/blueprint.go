package api

import (
	"context"

	"github.com/jinzhu/copier"
)

func (c *Client) ListBlueprints(ctx context.Context) ([]*BlueprintSummary, error) {
	out := make([]*BlueprintSummary, 0)

	_ = `# @genqlient
		fragment BlueprintSummary on Blueprint {
			id
			description
			displayName
			type
		}
		query getBlueprints($userId: ID!) {
			user(id: $userId) {
				blueprints {
					nodes {
						...BlueprintSummary
					}
				}
			}
		}`

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	res, err := getBlueprints(ctx, c.gql, user.Id)
	nodes := res.User.Blueprints.Nodes

	for _, n := range nodes {
		bps := &BlueprintSummary{}
		if err := copier.Copy(bps, n.BlueprintSummary); err != nil {
			return nil, err
		}
		out = append(out, bps)
	}

	return out, err
}
