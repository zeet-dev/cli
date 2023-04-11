package api

import (
	"context"
	"fmt"

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

func (c *Client) DeleteBlueprint(ctx context.Context, id uuid.UUID) error {
	_ = `# @genqlient
		mutation deleteBlueprint($id: UUID!) {
			deleteBlueprint(id: $id)
		}`

	res, err := deleteBlueprint(ctx, c.gql, id)
	if err != nil {
		return err
	}
	if !res.DeleteBlueprint {
		return fmt.Errorf("Failed to delete blueprint")
	}

	return nil
}

func (c *Client) ListBlueprints(ctx context.Context, pageInput PageInput) (*BlueprintConnection, error) {
	_ = `# @genqlient
		query getBlueprints($userId: ID!, $pageInput: PageInput!) {
			user(id: $userId) {
				# @genqlient(typename: "BlueprintConnection")
				blueprints(page: $pageInput) {
					totalCount
					# @genqlient(typename: "PageInfo")
					pageInfo {
						startCursor
						endCursor
						hasNextPage
						hasPreviousPage
					}
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

	res, err := getBlueprints(ctx, c.gql, user.Id, pageInput)
	bp := &res.User.Blueprints

	return bp, err
}
