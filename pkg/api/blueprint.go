package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// Shared fragment
var _ = `# @genqlient
		fragment BlueprintSummary on Blueprint {
			id
			slug
			displayName
			description
			type
			projectCount
			richInputSchema
			tags
		}`

func (c *Client) GetBlueprint(ctx context.Context, id uuid.UUID) (*BlueprintSummary, error) {
	out := &BlueprintSummary{}

	_ = `# @genqlient
		query getBlueprint($userID: ID!, $blueprintID: UUID!) {
			user(id: $userID) {
				blueprint(id: $blueprintID) {
					...BlueprintSummary
				}
			}
		}`

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}


	res, err := getBlueprint(ctx, c.gql, user.Id, id)
	if err := copier.Copy(out, res.User.Blueprint.BlueprintSummary); err != nil {
		return nil, err
	}

	return out, err
}

func (c *Client) ListBlueprints(ctx context.Context) ([]*BlueprintSummary, error) {
	out := make([]*BlueprintSummary, 0)

	_ = `# @genqlient
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
